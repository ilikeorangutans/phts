package model

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
	"github.com/nfnt/resize"
)

type Collection struct {
	db.CollectionRecord
	collectionRepo CollectionRepository
}

func (c Collection) AddPhoto(filename string, data []byte) error {
	if !c.IsPersisted() {
		return fmt.Errorf("Cannot add photos to unpersisted collection")
	}

	return c.collectionRepo.AddPhoto(c, filename, data)
}

type CollectionRepository interface {
	FindByID(id int64) (Collection, error)
	FindBySlug(slug string) (Collection, error)
	Save(Collection) (Collection, error)
	Create(string, string) Collection
	Recent(int) ([]Collection, error)
	AddPhoto(Collection, string, []byte) error
	AddRendition(db.PhotoRecord, db.RenditionRecord) (db.RenditionRecord, error)
	RecentPhotos(Collection, int) ([]Photo, error)
	DeletePhoto(Collection, Photo) error
}

func NewCollectionRepository(dbx *sqlx.DB, backend storage.Backend) CollectionRepository {
	return &collectionRepoImpl{
		db:          dbx,
		collections: db.NewCollectionDB(dbx),
		photos:      db.NewPhotoDB(dbx),
		renditions:  db.NewRenditionDB(dbx),
		backend:     backend,
		exifDB:      db.NewExifDB(dbx),
	}
}

type collectionRepoImpl struct {
	db          *sqlx.DB
	collections db.CollectionDB
	photos      db.PhotoDB
	renditions  db.RenditionDB
	backend     storage.Backend
	exifDB      db.ExifDB
}

func (r *collectionRepoImpl) AddRendition(photo db.PhotoRecord, rendition db.RenditionRecord) (db.RenditionRecord, error) {
	log.Printf("Adding rendition %v to photot %v", rendition, photo)

	tx, err := r.db.Beginx()
	if err != nil {
		return rendition, err
	}

	rendition, err = r.renditions.Save(rendition)
	if err != nil {
		tx.Rollback()
		return rendition, err
	}

	err = r.db.QueryRow("SELECT count(id) FROM renditions WHERE photo_id = $1", photo.ID).Scan(&photo.RenditionCount)
	if err != nil {
		tx.Rollback()
		return rendition, err
	}

	_, err = r.photos.Save(photo)
	if err != nil {
		tx.Rollback()
		return rendition, err
	}

	return rendition, tx.Commit()
}

func (r *collectionRepoImpl) DeletePhoto(col Collection, photo Photo) error {
	withTransaction(r.db, func() error {
		ids, err := r.photos.Delete(col.ID, photo.ID)
		if err != nil {
			return err
		}

		for _, id := range ids {
			if err := r.backend.Delete(id); err != nil {
				// TODO this sucks because deleting stuff cannot be rolled back.
				return err
			}
		}

		return nil
	})

	return nil
}

func (r *collectionRepoImpl) AddPhoto(collection Collection, filename string, data []byte) error {
	return withTransaction(r.db, func() error {
		var err error
		var takenAt *time.Time
		var tags ExifTags
		if tags, err = ExifTagsFromPhoto(data); err != nil {
			log.Printf("Could not extract EXIF from file %s", filename)
		} else {

			// TODO there's multiple date time tags, which one to use?
			if tag, err := tags.ByName("DateTimeOriginal"); err == nil {
				takenAt = tag.DateTime
			}
		}

		photo, err := r.photos.Save(db.PhotoRecord{
			CollectionID: collection.ID,
			Filename:     filename,
			TakenAt:      takenAt,
		})
		if err != nil {
			return err
		}

		for _, tag := range tags {
			log.Printf("Saving EXIF %s", tag.ExifRecord)
			_, err = r.exifDB.Save(photo.ID, tag.ExifRecord)
			if err != nil {
				return err
			}
		}

		rendition, err := db.NewRenditionRecord(photo, filename, data)
		if err != nil {
			return err
		}
		rendition.Original = true
		rendition, err = r.renditions.Save(rendition)
		if err != nil {
			return err
		}

		err = r.backend.Store(rendition.ID, data)
		if err != nil {
			return err
		}

		log.Printf("Created photo %d with rendition %d", photo.ID, rendition.ID)

		_, err = r.Save(collection)

		configs, err := r.renditions.ApplicableConfigs(collection.ID)
		if err != nil {
			return err
		}

		for _, config := range configs {
			log.Printf("Config %v", config)
			go makeThumbnail(r, r.backend, photo, filename, data, uint(config.Width))
		}
		return err
	})
}

func makeThumbnail(r CollectionRepository, backend storage.Backend, photo db.PhotoRecord, filename string, data []byte, maxSize uint) {
	log.Printf("Creating thumbnail for photo %v, maxSize %d", photo, maxSize)
	rawJpeg, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("Could not resize file %s for photo %v: %s", filename, photo, err)
		return
	}

	resized := resize.Resize(maxSize, 0, rawJpeg, resize.Lanczos3)

	record := db.RenditionRecord{
		Timestamps: db.JustCreated(),
		PhotoID:    photo.ID,
		Original:   false,
		Width:      uint(resized.Bounds().Dx()),
		Height:     uint(resized.Bounds().Dy()),
		Format:     "image/jpeg",
	}

	record, err = r.AddRendition(photo, record)
	if err != nil {
		log.Printf("Could not resize file %s for photo %v: %s", filename, photo, err)
		return
	}

	var b = &bytes.Buffer{}
	// TODO use quality settings from rendition configuration
	err = jpeg.Encode(b, resized, &jpeg.Options{Quality: 95})
	if err != nil {
		log.Printf("Could not resize file %s for photo %v: %s", filename, photo, err)
		return
	}

	backend.Store(record.ID, b.Bytes())
}

func withTransaction(db *sqlx.DB, f func() error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = f()
	if err != nil {
		rollbackErr := tx.Rollback()
		log.Println(rollbackErr)
		return err
	}

	if err = tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()
		log.Println(rollbackErr)
		return err
	}
	return nil
}

func (r *collectionRepoImpl) RecentPhotos(collection Collection, count int) ([]Photo, error) {
	photos, err := r.photos.List(collection.ID, 0, "updated_at", count)
	if err != nil {
		return nil, err
	}

	photo_ids := []int64{}
	for _, photo := range photos {
		photo_ids = append(photo_ids, photo.ID)
	}

	// TODO: hardcoded photo size here, should find the thumbail
	rends, err := r.renditions.FindBySize(photo_ids, 345, 0)
	if err != nil {
		return nil, err
	}

	var result []Photo

	for _, photoRecord := range photos {
		renditions := []Rendition{}
		if rend, ok := rends[photoRecord.ID]; ok {
			renditions = append(renditions, Rendition{rend})
		}
		photo := Photo{
			PhotoRecord: photoRecord,
			Renditions:  renditions,
			Collection:  collection,
		}

		result = append(result, photo)
	}

	return result, nil
}

func (r *collectionRepoImpl) Recent(count int) ([]Collection, error) {
	records, err := r.collections.List(count, 0, "updated_at")
	if err != nil {
		return nil, err
	}

	result := []Collection{}
	for _, record := range records {
		col := r.newCollection()
		col.CollectionRecord = record
		result = append(result, col)
	}

	return result, nil
}

func (r *collectionRepoImpl) FindByID(id int64) (Collection, error) {
	if record, err := r.collections.FindByID(id); err != nil {
		return Collection{}, err
	} else {
		return Collection{
			CollectionRecord: record,
			collectionRepo:   r,
		}, nil
	}
}

func (r *collectionRepoImpl) FindBySlug(slug string) (Collection, error) {
	if record, err := r.collections.FindBySlug(slug); err != nil {
		return Collection{}, err
	} else {
		return Collection{
			CollectionRecord: record,
			collectionRepo:   r,
		}, nil
	}
}

func (r *collectionRepoImpl) Save(collection Collection) (Collection, error) {
	record, err := r.collections.Save(collection.CollectionRecord)
	collection.CollectionRecord = record
	return collection, err
}

func (r *collectionRepoImpl) newCollection() Collection {
	return Collection{
		CollectionRecord: db.CollectionRecord{},
		collectionRepo:   r,
	}
}

func (r *collectionRepoImpl) Create(name string, slug string) Collection {
	result := r.newCollection()
	result.Name = name
	result.Slug = slug
	return result
}

func DBFromRequest(r *http.Request) *sqlx.DB {
	db, ok := r.Context().Value("database").(*sqlx.DB)
	if !ok {
		log.Fatal("Could not get database from request, wrong type")
	}

	return db
}

func CollectionRepoFromRequest(r *http.Request) CollectionRepository {
	backend, ok := r.Context().Value("backend").(storage.Backend)
	if !ok {
		log.Fatal("Could not get backend from request, wrong type")
	}

	db := DBFromRequest(r)
	return NewCollectionRepository(db, backend)
}

type PhotoRepository interface {
	FindByID(collection Collection, photoID int64) (Photo, error)
	List(collection Collection, afterID int64, count int, sortBy string) ([]Photo, int64, error)
}

type photoRepoImpl struct {
	backend    storage.Backend
	db         *sqlx.DB
	photos     db.PhotoDB
	renditions db.RenditionDB
	exifDB     db.ExifDB
}

func (r *photoRepoImpl) FindByID(collection Collection, photoID int64) (Photo, error) {
	record, err := r.photos.FindByID(collection.ID, photoID)
	if err != nil {
		return Photo{}, err
	}

	renditions, err := r.renditions.FindAllForPhoto(record.ID)
	if err != nil {
		return Photo{}, err
	}

	exifTags, err := r.exifDB.AllForPhoto(photoID)
	if err != nil {
		return Photo{}, err
	}

	photo := Photo{
		PhotoRecord: record,
		Renditions:  []Rendition{},
		Exif:        []ExifTag{},
	}

	for _, rendition := range renditions {
		photo.Renditions = append(photo.Renditions, Rendition{rendition})
	}

	for _, tag := range exifTags {
		photo.Exif = append(photo.Exif, ExifTag{tag})
	}

	return photo, err
}

func (r *photoRepoImpl) List(collection Collection, afterID int64, count int, sortBy string) (photos []Photo, lastID int64, err error) {
	records, err := r.photos.List(collection.ID, afterID, sortBy, count)
	if err != nil {
		return nil, 0, err
	}

	config, err := r.renditions.FindConfig(collection.ID, "admin thumbnails")
	if err != nil {
		return nil, 0, err
	}

	photoIDs := []int64{}
	for _, p := range records {
		photoIDs = append(photoIDs, p.ID)
	}

	renditions, err := r.renditions.FindBySize(photoIDs, config.Width, 0)
	if err != nil {
		return nil, 0, err
	}

	result := []Photo{}

	for _, p := range records {
		result = append(result, Photo{
			PhotoRecord: p,
			Renditions:  []Rendition{Rendition{renditions[p.ID]}},
			Collection:  collection,
		})
		lastID = p.ID
	}

	return result, lastID, nil
}

func PhotoRepoFromRequest(r *http.Request) PhotoRepository {
	backend, ok := r.Context().Value("backend").(storage.Backend)
	if !ok {
		log.Fatal("Could not get backend from request, wrong type")
	}

	dbx := DBFromRequest(r)
	return &photoRepoImpl{
		db:         dbx,
		backend:    backend,
		photos:     db.NewPhotoDB(dbx),
		renditions: db.NewRenditionDB(dbx),
		exifDB:     db.NewExifDB(dbx),
	}
}
