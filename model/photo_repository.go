package model

import (
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
)

type PhotoRepository interface {
	FindByID(collection Collection, photoID int64) (Photo, error)
	List(collection Collection, paginator db.Paginator) ([]Photo, db.Paginator, error)
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
