package model

import (
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/storage"
)

// TODO this does not work as intended
// we're still using the actual db reference but what we ought to do is use the transaction instead
func withTransaction(db db.DB, f func() error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = f()
	if err != nil {
		log.Printf("Rolling back transaction due to error: %s", err.Error())
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

func DBFromRequest(r *http.Request) db.DB {
	db, ok := r.Context().Value("database").(db.DB)
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

func PhotoRepoFromRequest(r *http.Request) PhotoRepository {
	backend, ok := r.Context().Value("backend").(storage.Backend)
	if !ok {
		log.Fatal("Could not get backend from request, wrong type")
	}

	dbx := DBFromRequest(r)
	return &photoRepoImpl{
		db:         dbx,
		backend:    backend,
		photos:     db.NewPhotoDB(dbx),
		renditions: db.NewRenditionDB(dbx),
		exifDB:     db.NewExifDB(dbx),
	}
}
