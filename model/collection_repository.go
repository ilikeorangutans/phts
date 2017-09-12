package model

import (
	"log"
	"time"

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
	AddPhoto(Collection, string, []byte) error
	AddRendition(db.PhotoRecord, db.RenditionRecord) (db.RenditionRecord, error)
	RecentPhotos(Collection, int) ([]Photo, error)
	DeletePhoto(Collection, Photo) error
	Delete(Collection) error
}

func NewCollectionRepository(dbx db.DB, backend storage.Backend) CollectionRepository {
	return &collectionRepoImpl{
		db:               dbx,
		collections:      db.NewCollectionDB(dbx),
		photos:           db.NewPhotoDB(dbx),
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
	renditions       db.RenditionDB
	renditionConfigs db.RenditionConfigurationDB
	backend          storage.Backend
	exifDB           db.ExifDB
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
		//ids, err := r.photos.Delete(col.ID, photo.ID)
		//if err != nil {
		//return err
		//}
		//
		//for _, id := range ids {
		//if err := r.backend.Delete(id); err != nil {
		//// TODO this sucks because deleting stuff cannot be rolled back.
		//return err
		//}
		//}

		return nil
	})

	return nil
}

func (r *collectionRepoImpl) AddPhoto(collection Collection, filename string, data []byte) error {
	// TODO a lot of the functionality in here should be in other repos.
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

		if err = r.createOriginalRendition(photo, filename, data); err != nil {
			return err
		}

		configs, err := r.renditionConfigs.FindForCollection(collection.ID)
		if err != nil {
			return err
		}

		for _, renditionConfiguration := range configs {
			// TODO Ideally we'd do this in the background (or push this to a queue)
			makeThumbnail(r, r.backend, photo, filename, data, uint(renditionConfiguration.Width))
		}

		_, err = r.Save(collection)

		return err
	})
}

func (r *collectionRepoImpl) createOriginalRendition(photo db.PhotoRecord, filename string, data []byte) error {
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
	result.collectionRepo = r
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
