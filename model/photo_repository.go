package model

import (
	"log"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/storage"
)

type PhotoRepository interface {
	FindByID(collection Collection, photoID int64) (Photo, error)
	List(collection Collection, paginator db.Paginator) ([]Photo, db.Paginator, error)
	AddRendition(Photo, db.RenditionRecord) (db.RenditionRecord, error)
	// Create adds a new photo to the given collection.
	Create(Collection, string, []byte) (Photo, error)
}

func NewPhotoRepository(dbx db.DB, backend storage.Backend) PhotoRepository {
	return &photoRepoImpl{
		backend:          backend,
		db:               dbx,
		photos:           db.NewPhotoDB(dbx),
		renditions:       db.NewRenditionDB(dbx),
		exifDB:           db.NewExifDB(dbx),
		renditionConfigs: db.NewRenditionConfigurationDB(dbx),
	}
}

type photoRepoImpl struct {
	backend          storage.Backend
	db               db.DB
	photos           db.PhotoDB
	renditions       db.RenditionDB
	exifDB           db.ExifDB
	renditionConfigs db.RenditionConfigurationDB
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

func (r *photoRepoImpl) List(collection Collection, paginator db.Paginator) ([]Photo, db.Paginator, error) {
	records, err := r.photos.List(collection.ID, paginator)
	if err != nil {
		return nil, paginator, err
	}

	//config, err := r.renditions.FindConfig(collection.ID, "admin thumbnails")
	//if err != nil {
	//return nil, paginator, err
	//}
	config := db.RenditionConfigurationRecord{}

	photoIDs := []int64{}
	for _, p := range records {
		photoIDs = append(photoIDs, p.ID)
	}

	if len(photoIDs) == 0 {
		return nil, paginator, nil
	}

	renditions, err := r.renditions.FindBySize(photoIDs, config.Width, 0)
	if err != nil {
		return nil, paginator, err
	}

	result := []Photo{}

	for _, p := range records {
		result = append(result, Photo{
			PhotoRecord: p,
			Renditions:  []Rendition{Rendition{renditions[p.ID]}},
			Collection:  collection,
		})
		paginator.PrevID = p.ID
		// TODO hardcoded value here, should use the column configured in the paginator
		paginator.PrevTimestamp = &p.UpdatedAt
	}

	return result, paginator, nil
}

func (r *photoRepoImpl) AddRendition(photo Photo, rendition db.RenditionRecord) (db.RenditionRecord, error) {
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

	_, err = r.photos.Save(photo.PhotoRecord)
	if err != nil {
		tx.Rollback()
		return rendition, err
	}

	return rendition, tx.Commit()
}

func (r *photoRepoImpl) Create(collection Collection, filename string, data []byte) (Photo, error) {
	var photo Photo
	err := withTransaction(r.db, func() error {
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

		photoRecord, err := r.photos.Save(db.PhotoRecord{
			CollectionID: collection.ID,
			Filename:     filename,
			TakenAt:      takenAt,
		})
		if err != nil {
			return err
		}

		for _, tag := range tags {
			log.Printf("Saving EXIF %s", tag.ExifRecord)
			_, err = r.exifDB.Save(photoRecord.ID, tag.ExifRecord)
			if err != nil {
				return err
			}
		}

		if err = r.createOriginalRendition(photoRecord, filename, data); err != nil {
			return err
		}

		configs, err := r.renditionConfigs.FindForCollection(collection.ID)
		if err != nil {
			return err
		}

		photo = Photo{
			Collection:  collection,
			PhotoRecord: photoRecord,
			Exif:        tags,
		}

		for _, renditionConfiguration := range configs {
			// TODO Ideally we'd do this in the background (or push this to a queue)
			makeThumbnail(r, r.backend, photo, filename, data, uint(renditionConfiguration.Width))
		}

		return err
	})
	return photo, err
}

func (r *photoRepoImpl) createOriginalRendition(photo db.PhotoRecord, filename string, data []byte) error {
	rendition, err := r.renditions.Create(photo, filename, data)
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

	return nil
}
