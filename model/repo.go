package model

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"

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

	log.Printf("Adding new photo %s with %d bytes to collection %s", filename, len(data), c.Name)

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
	RecentPhotos(Collection) ([]Photo, error)
}

func NewCollectionRepository(dbx *sqlx.DB, backend storage.Backend) CollectionRepository {
	return &collectionRepoImpl{
		db:          dbx,
		collections: db.NewCollectionDB(dbx),
		photos:      db.NewPhotoDB(dbx),
		renditions:  db.NewRenditionDB(dbx),
		backend:     backend,
	}
}

type collectionRepoImpl struct {
	db          *sqlx.DB
	collections db.CollectionDB
	photos      db.PhotoDB
	renditions  db.RenditionDB
	backend     storage.Backend
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

func (r *collectionRepoImpl) AddPhoto(collection Collection, filename string, data []byte) error {
	return withTransaction(r.db, func() error {
		photo, err := r.photos.Save(db.PhotoRecord{
			CollectionID: collection.ID,
		})
		if err != nil {
			return err
		}

		// TODO here we'd extract EXIF
		rendition, err := db.NewRenditionRecord(photo, filename, data)
		if err != nil {
			return err
		}
		rendition, err = r.renditions.Save(rendition)
		if err != nil {
			return err
		}

		err = r.backend.Store(rendition.ID, data)
		if err != nil {
			return err
		}

		log.Printf("Created photo %d with rendition %d", photo.ID, rendition.ID)

		collection.PhotoCount += 1 // TODO: better to actually count
		_, err = r.Save(collection)

		go makeThumbnail(r, r.backend, photo, filename, data, 256)
		go makeThumbnail(r, r.backend, photo, filename, data, 1024)
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

func (r *collectionRepoImpl) RecentPhotos(collection Collection) ([]Photo, error) {

	photos, err := r.photos.List(collection.ID, 0, "updated_at asc", 10)
	if err != nil {
		return nil, err
	}

	//photosAndRenditions, err := r.photos.ListWithRenditions(collection.ID, 10)
	//if err != nil {
	//return nil, err
	//}

	photo_ids := []int64{}
	for _, photo := range photos {
		photo_ids = append(photo_ids, photo.ID)
	}

	rends, err := r.renditions.FindBySize(photo_ids, 256, 0)
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
	FindByID(collectionID int64, photoID int64) (Photo, error)
	Save(photo Photo) (Photo, error)
}

type PhotoSQLRepository struct {
	backend storage.Backend
	db      *sqlx.DB
}

func (r *PhotoSQLRepository) FindByID(collectionID, photoID int64) (Photo, error) {
	var photo Photo
	err := r.db.QueryRowx("SELECT * FROM photos WHERE collection_id=$1 AND photo_id=$2", collectionID, photoID).StructScan(&photo)
	return photo, err
}
