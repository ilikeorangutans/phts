package model

import (
	"log"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/ilikeorangutans/phts/pkg/metadata"
	"github.com/ilikeorangutans/phts/storage"
)

type PhotoRepository interface {
	FindByID(collection *db.Collection, photoID int64) (Photo, error)
	List(collection *db.Collection, paginator database.Paginator, configs []RenditionConfiguration) ([]Photo, database.Paginator, error)
	ListAlbum(collection *db.Collection, album Album, paginator database.Paginator, configs []RenditionConfiguration) ([]Photo, database.Paginator, error)
	Delete(collection *db.Collection, photo Photo) error
	// Create adds a new photo to the given collection.
	Create(*db.Collection, string, []byte) (Photo, error)
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

func (r *photoRepoImpl) FindByID(collection *db.Collection, photoID int64) (Photo, error) {
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
		Exif:        []metadata.ExifTag{},
	}

	for _, rendition := range renditions {
		photo.Renditions = append(photo.Renditions, Rendition{rendition, nil})
	}

	for _, tag := range exifTags {
		photo.Exif = append(photo.Exif, metadata.ExifTag{tag})
	}

	return photo, err
}

func (r *photoRepoImpl) List(collection *db.Collection, paginator database.Paginator, renditionConfigs []RenditionConfiguration) ([]Photo, database.Paginator, error) {
	records, err := r.photos.List(collection.ID, paginator)
	if err != nil {
		return nil, paginator, err
	}

	photoIDs := []int64{}
	for _, p := range records {
		photoIDs = append(photoIDs, p.ID)
	}

	if len(photoIDs) == 0 {
		return []Photo{}, paginator, nil
	}

	renditionConfigIDs := []int64{}
	for _, rc := range renditionConfigs {
		renditionConfigIDs = append(renditionConfigIDs, rc.ID)
	}

	if len(renditionConfigIDs) == 0 {
		return []Photo{}, paginator, nil
	}

	renditionRecords, err := r.renditions.FindByRenditionConfigurations(photoIDs, renditionConfigIDs)
	if err != nil {
		return nil, paginator, err
	}

	renditions := make(map[int64]Renditions)
	for photoID, records := range renditionRecords {
		for _, record := range records {
			rendition := Rendition{record, nil}
			renditions[photoID] = append(renditions[photoID], rendition)
		}
	}

	result := []Photo{}

	for _, p := range records {
		result = append(result, NewPhotoFromRecord(p, collection, renditions[p.ID]))
		paginator.PrevID = p.ID
		// TODO hardcoded value here, should use the column configured in the paginator
		paginator.PrevTimestamp = &p.UpdatedAt
	}

	return result, paginator, nil
}

func (r *photoRepoImpl) ListAlbum(collection *db.Collection, album Album, paginator database.Paginator, renditionConfigs []RenditionConfiguration) ([]Photo, database.Paginator, error) {
	records, err := r.photos.ListAlbum(collection.ID, album.ID, paginator)
	if err != nil {
		return nil, paginator, err
	}

	photoIDs := []int64{}
	for _, p := range records {
		photoIDs = append(photoIDs, p.ID)
	}

	if len(photoIDs) == 0 {
		return []Photo{}, paginator, nil
	}

	renditionConfigIDs := []int64{}
	for _, rc := range renditionConfigs {
		renditionConfigIDs = append(renditionConfigIDs, rc.ID)
	}

	if len(renditionConfigIDs) == 0 {
		return []Photo{}, paginator, nil
	}

	renditionRecords, err := r.renditions.FindByRenditionConfigurations(photoIDs, renditionConfigIDs)
	if err != nil {
		return nil, paginator, err
	}

	renditions := make(map[int64]Renditions)
	for photoID, records := range renditionRecords {
		for _, record := range records {
			rendition := Rendition{record, nil}
			renditions[photoID] = append(renditions[photoID], rendition)
		}
	}

	result := []Photo{}

	for _, p := range records {
		result = append(result, NewPhotoFromRecord(p, collection, renditions[p.ID]))
		paginator.PrevID = p.ID
		// TODO hardcoded value here, should use the column configured in the paginator
		paginator.PrevTimestamp = &p.UpdatedAt
	}

	return result, paginator, nil

}

func (r *photoRepoImpl) Delete(collection *db.Collection, photo Photo) error {
	return r.photos.Delete(collection.ID, photo.ID)
}

func (r *photoRepoImpl) Create(collection *db.Collection, filename string, data []byte) (Photo, error) {
	var photo Photo
	err := withTransaction(r.db, func() error {
		var err error
		var takenAt *time.Time
		var tags metadata.ExifTags
		// TODO removing this soon
		//		if tags, err = metadata.ExifTagsFromPhoto(data); err != nil {
		//			log.Printf("Could not extract EXIF from file %s", filename)
		//		} else {
		//
		//			takenAtFields := []string{"DateTime", "DateTimeOriginal"}
		//			for _, field := range takenAtFields {
		//				if tag, err := tags.ByName(field); err == nil {
		//					takenAt = tag.DateTime
		//					break
		//				}
		//			}
		//		}
		//
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
				log.Printf("error saving exif tag %s: %s", tag.Tag, err.Error())
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

		orientation := metadata.Horizontal
		if orientationTag, err := tags.ByName("Orientation"); err == nil {
			orientation = metadata.ExifOrientationFromTag(orientationTag)
			log.Println(orientation)
		}

		renditions, err := configs.Process(filename, data, orientation)
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
			//Collection:  collection,
			PhotoRecord: photoRecord,
			Exif:        tags,
		}

		return err
	})
	return photo, err
}
