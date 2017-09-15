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
		photo.Renditions = append(photo.Renditions, Rendition{rendition, nil})
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
			Renditions:  []Rendition{Rendition{renditions[p.ID], nil}},
			Collection:  collection,
		})
		paginator.PrevID = p.ID
		// TODO hardcoded value here, should use the column configured in the paginator
		paginator.PrevTimestamp = &p.UpdatedAt
	}

	return result, paginator, nil
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
			_, err = r.exifDB.Save(photoRecord.ID, tag.ExifRecord)
			if err != nil {
				return err
			}
		}

		configRecords, err := r.renditionConfigs.FindForCollection(collection.ID)
		if err != nil {
			return err
		}
		var configs RenditionConfigurations
		for _, record := range configRecords {
			configs = append(configs, RenditionConfiguration{record})
		}

		renditions, err := configs.Process(filename, data)
		if err != nil {
			return err
		}

		for _, rendition := range renditions {
			rendition.PhotoID = photoRecord.ID
			renditionRecord, err := r.renditions.Save(rendition.RenditionRecord)
			if err != nil {
				return err
			}

			if err := r.backend.Store(renditionRecord.ID, rendition.data); err != nil {
				return err
			}
		}

		photoRecord, err = r.photos.Save(photoRecord)
		photo = Photo{
			Collection:  collection,
			PhotoRecord: photoRecord,
			Exif:        tags,
		}

		return err
	})
	return photo, err
}
