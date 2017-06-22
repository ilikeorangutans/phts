package model

import (
	"log"
	"net/http"
	"time"

	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
)

type CollectionRepository interface {
	FindByID(id uint) (Collection, error)
	FindBySlug(slug string) (Collection, error)
	Save(Collection) (Collection, error)
	Recent(int) ([]Collection, error)
}

type CollectionSQLRepository struct {
	create          func() Collection
	db              *sqlx.DB
	photoRepository PhotoRepository
}

func (r *CollectionSQLRepository) FindBySlug(slug string) (Collection, error) {
	var result Collection
	err := r.db.QueryRowx("SELECT * FROM collections WHERE slug=$1", slug).StructScan(&result)
	return result, err
}

func (r *CollectionSQLRepository) FindByID(id uint) (Collection, error) {
	var result Collection
	err := r.db.QueryRow("SELECT * FROM collections WHERE id=$1", id).Scan(&result)
	return result, err
}

func (r *CollectionSQLRepository) Recent(n int) ([]Collection, error) {

	rows, err := r.db.Queryx("SELECT * FROM collections ORDER BY updated_at DESC LIMIT $1", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Collection{}
	for rows.Next() {
		col := Collection{}
		err = rows.StructScan(&col)
		if err != nil {
			return nil, err
		}
		result = append(result, col)
	}

	return result, nil
}

// TODO: runWithTransaction(func) would be super useful

func (r *CollectionSQLRepository) Save(col Collection) (Collection, error) {
	if col.ID == 0 {
		tx, err := r.db.Beginx()
		if err != nil {
			return col, err
		}
		now := time.Now()
		err = r.db.QueryRow("UPDATE collections (name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", col.Name, col.Slug, now, now).Scan(&col.ID)

		for _, photo := range col.Photos {
			if _, err = r.photoRepository.Save(photo); err != nil {
				tx.Rollback()
				return col, err
			}
		}

		return col, tx.Commit()
	} else {
		tx, err := r.db.Beginx()
		if err != nil {
			return col, err
		}

		now := time.Now()
		err = r.db.QueryRow("INSERT INTO collections (name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", col.Name, col.Slug, now, now).Scan(&col.ID)
		if err != nil {
			tx.Rollback()
			return col, err
		}

		return col, tx.Commit()
	}
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

	photoRepository := &PhotoSQLRepository{
		db:      db,
		backend: backend,
	}

	return &CollectionSQLRepository{
		db: db,
		create: func() Collection {
			return Collection{}
		},
		photoRepository: photoRepository,
	}
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

func (r *PhotoSQLRepository) Save(photo Photo) (Photo, error) {
	if photo.ID == 0 {
		return r.saveNew(photo)
	} else {
		return r.updateExisting(photo)
	}
}

func (r *PhotoSQLRepository) updateExisting(photo Photo) (Photo, error) {
	// TODO check for changes etc
	return photo, nil
}

func (r *PhotoSQLRepository) saveNew(photo Photo) (Photo, error) {
	log.Println("Saving new photo")
	tx, err := r.db.Beginx()
	if err != nil {
		log.Fatal(err)
	}
	err = r.db.QueryRow("INSERT INTO photos (collection_id, updated_at) VALUES ($1, $2) RETURNING id", photo.CollectionID, time.Now()).Scan(&photo.ID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			log.Printf("Error rolling back transaction %s", err)
		}
		log.Fatal(err)
	}

	for _, rendition := range photo.Renditions {
		err = r.db.QueryRow("INSERT INTO renditions (photo_id, original, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", photo.ID, rendition.Original, rendition.CreatedAt, rendition.UpdatedAt).Scan(&rendition.ID)
		if err != nil {
			log.Fatal(err)
		}

		err = r.backend.Store(rendition.ID, rendition.Data)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return photo, nil
}
