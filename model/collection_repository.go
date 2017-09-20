package model

import (
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/storage"
)

type CollectionRepository interface {
	FindByID(id int64) (Collection, error)
	FindBySlug(slug string) (Collection, error)
	// Save saves or updates a given collection
	Save(Collection) (Collection, error)
	// Create a new instance of Collection.
	Create(name, slug string) Collection
	Recent(int) ([]Collection, error)
	AddPhoto(Collection, string, []byte) (Photo, error)
	RecentPhotos(Collection, int) ([]Photo, error)
	DeletePhoto(Collection, Photo) error
	Delete(Collection) error
	ApplicableRenditionConfigurations(Collection) ([]RenditionConfiguration, error)
}

func NewCollectionRepository(dbx db.DB, backend storage.Backend) CollectionRepository {
	return &collectionRepoImpl{
		db:               dbx,
		collections:      db.NewCollectionDB(dbx),
		photos:           db.NewPhotoDB(dbx),
		photoRepo:        NewPhotoRepository(dbx, backend),
		renditions:       db.NewRenditionDB(dbx),
		renditionConfigs: db.NewRenditionConfigurationDB(dbx),
		backend:          backend,
		exifDB:           db.NewExifDB(dbx),
	}
}

type collectionRepoImpl struct {
	db               db.DB
	collections      db.CollectionDB
	photos           db.PhotoDB
	photoRepo        PhotoRepository
	renditions       db.RenditionDB
	renditionConfigs db.RenditionConfigurationDB
	backend          storage.Backend
	exifDB           db.ExifDB
}

func (r *collectionRepoImpl) ApplicableRenditionConfigurations(col Collection) ([]RenditionConfiguration, error) {
	records, err := r.renditionConfigs.FindForCollection(col.ID)
	if err != nil {
		return nil, err
	}

	var result []RenditionConfiguration
	for _, record := range records {
		result = append(result, RenditionConfiguration{record})
	}
	return result, nil
}

func (r *collectionRepoImpl) DeletePhoto(col Collection, photo Photo) error {
	// TODO withTransaction does not work as intended. Fix this
	withTransaction(r.db, func() error {
		renditionIDs, err := r.renditions.DeleteForPhoto(photo.ID)
		if err != nil {
			return err
		}

		err = r.photos.Delete(col.ID, photo.ID)
		if err != nil {
			log.Printf("Delete failed %s", err.Error())
			return err
		}

		for _, id := range renditionIDs {
			if err := r.backend.Delete(id); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func (r *collectionRepoImpl) AddPhoto(collection Collection, filename string, data []byte) (Photo, error) {
	if photo, err := r.photoRepo.Create(collection, filename, data); err != nil {
		return photo, err
	} else {
		_, err = r.Save(collection)
		return photo, err
	}
}

func (r *collectionRepoImpl) RecentPhotos(collection Collection, count int) ([]Photo, error) {
	paginator := db.NewPaginator()
	paginator.Count = uint(count)

	photos, err := r.photos.List(collection.ID, paginator)
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
			renditions = append(renditions, Rendition{rend, nil})
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
		}, nil
	}
}

func (r *collectionRepoImpl) FindBySlug(slug string) (Collection, error) {
	if record, err := r.collections.FindBySlug(slug); err != nil {
		return Collection{}, err
	} else {
		return Collection{
			CollectionRecord: record,
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
	}
}

func (r *collectionRepoImpl) Create(name string, slug string) Collection {
	result := r.newCollection()
	result.Name = name
	result.Slug = slug
	return result
}

func (r *collectionRepoImpl) Delete(collection Collection) error {
	if !collection.IsPersisted() {
		// TODO should this be an error?
		return nil
	}

	// TODO do other cleanup work
	// TODO Get rendition IDs and delete from storage

	return r.collections.Delete(collection.CollectionRecord.ID)
}
