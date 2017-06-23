package model

import (
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
)

type Collection struct {
	db.CollectionRecord
	collectionRepo CollectionRepository
}

type CollectionRepository interface {
	FindByID(id int64) (Collection, error)
	FindBySlug(slug string) (Collection, error)
	Save(Collection) (Collection, error)
	Create(string, string) Collection
	Recent(int) ([]Collection, error)
}

func NewCollectionRepository(dbx *sqlx.DB) CollectionRepository {
	return &collectionRepoImpl{
		db:          dbx,
		collections: db.NewCollectionDB(dbx),
	}
}

type collectionRepoImpl struct {
	db          *sqlx.DB
	collections db.CollectionDB
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
	_, ok := r.Context().Value("backend").(storage.Backend)
	if !ok {
		log.Fatal("Could not get backend from request, wrong type")
	}

	db := DBFromRequest(r)
	return NewCollectionRepository(db)
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
