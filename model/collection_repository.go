package model

import (
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/storage"
)

type CollectionFinder interface {
	FindByID(id int64) (db.CollectionRecord, error)
	FindBySlug(slug string) (db.CollectionRecord, error)
}

// CollectionRepository allows access to collections.
type CollectionRepository interface {
	CollectionFinder
	// Save saves or updates a given collection
	Save(db.CollectionRecord) (db.CollectionRecord, error)
	// Create a new instance of db.CollectionRecord.
	Create(name, slug string) db.CollectionRecord
	Recent(int) ([]db.CollectionRecord, error)
	AddPhoto(db.CollectionRecord, string, []byte) (Photo, error)
	DeletePhoto(db.CollectionRecord, Photo) error
	Delete(db.CollectionRecord) error
	ApplicableRenditionConfigurations(db.CollectionRecord) (RenditionConfigurations, error)
	// Remove the given configuration from the collection; this will delete all associated renditions.
	RemoveRenditionConfiguration(db.CollectionRecord, RenditionConfiguration) error
}

// NewUserCollectionRepository returns a CollectionRepository for a specific user. All operations are scoped
// to the user passed in.
func NewUserCollectionRepository(dbx db.DB, backend storage.Backend, user *db.UserRecord) CollectionRepository {
	return &userCollectionRepoImpl{
		db:               dbx,
		collections:      db.NewCollectionDB(dbx),
		photos:           db.NewPhotoDB(dbx),
		photoRepo:        NewPhotoRepository(dbx, backend),
		renditions:       db.NewRenditionDB(dbx),
		renditionConfigs: db.NewRenditionConfigurationDB(dbx),
		backend:          backend,
		exifDB:           db.NewExifDB(dbx),
		user:             user,
	}
}

type userCollectionRepoImpl struct {
	db               db.DB
	collections      db.CollectionDB
	photos           db.PhotoDB
	photoRepo        PhotoRepository
	renditions       db.RenditionDB
	renditionConfigs db.RenditionConfigurationDB
	backend          storage.Backend
	exifDB           db.ExifDB
	user             *db.UserRecord
}

func (r *userCollectionRepoImpl) canAccess(col db.CollectionRecord) bool {
	// TODO we'd call this method before accessing collections
	return r.collections.CanAccess(r.user.ID, col.ID)
}

func (r *userCollectionRepoImpl) ApplicableRenditionConfigurations(col db.CollectionRecord) (RenditionConfigurations, error) {
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

func (r *userCollectionRepoImpl) DeletePhoto(col db.CollectionRecord, photo Photo) error {
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

func (r *userCollectionRepoImpl) AddPhoto(collection db.CollectionRecord, filename string, data []byte) (Photo, error) {
	photo, err := r.photoRepo.Create(collection, filename, data)
	if err != nil {
		return photo, err
	}

	_, err = r.Save(collection)
	return photo, err
}

func (r *userCollectionRepoImpl) Recent(count int) ([]db.CollectionRecord, error) {
	records, err := r.collections.List(r.user.ID, count, 0, "updated_at")
	if err != nil {
		return nil, err
	}

	result := []db.CollectionRecord{}
	for _, record := range records {
		result = append(result, record)
	}

	return result, nil
}

func (r *userCollectionRepoImpl) FindByID(id int64) (db.CollectionRecord, error) {
	if record, err := r.collections.FindByID(id); err != nil {
		return db.CollectionRecord{}, err
	} else {
		return record, nil
	}
}

func (r *userCollectionRepoImpl) FindBySlug(slug string) (db.CollectionRecord, error) {
	if record, err := r.collections.FindBySlug(slug); err != nil {
		return db.CollectionRecord{}, err
	} else {
		return record, nil
	}
}

func (r *userCollectionRepoImpl) Save(collection db.CollectionRecord) (db.CollectionRecord, error) {
	newRecord := !collection.IsPersisted()
	record, err := r.collections.Save(collection)
	if err != nil {
		return collection, err
	}
	if newRecord {
		err = r.collections.Assign(r.user.ID, record.ID)
	}
	return record, err
}

func (r *userCollectionRepoImpl) Create(name string, slug string) db.CollectionRecord {
	result := db.CollectionRecord{}
	result.Name = name
	result.Slug = slug
	return result
}

func (r *userCollectionRepoImpl) RemoveRenditionConfiguration(collection db.CollectionRecord, config RenditionConfiguration) error {
	// TODO implement me
	return nil
}

func (r *userCollectionRepoImpl) Delete(collection db.CollectionRecord) error {
	if !collection.IsPersisted() {
		// TODO should this be an error?
		return nil
	}

	// TODO do other cleanup work
	// TODO Get rendition IDs and delete from storage

	return r.collections.Delete(collection.ID)
}
