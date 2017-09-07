package model

import (
	"bytes"
	"image/jpeg"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
	"github.com/nfnt/resize"
)

func makeThumbnail(r CollectionRepository, backend storage.Backend, photo db.PhotoRecord, filename string, data []byte, maxSize uint) {
	log.Printf("Creating thumbnail for photo %v, maxSize %d", photo, maxSize)
	rawJpeg, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("Could not resize file %s for photo %v: %s", filename, photo, err)
		return
	}

	resized := resize.Resize(maxSize, 0, rawJpeg, resize.Lanczos3)

	record := db.RenditionRecord{
		Timestamps: db.JustCreated(),
		PhotoID:    photo.ID,
		Original:   false,
		Width:      uint(resized.Bounds().Dx()),
		Height:     uint(resized.Bounds().Dy()),
		Format:     "image/jpeg",
	}

	record, err = r.AddRendition(photo, record)
	if err != nil {
		log.Printf("Could not resize file %s for photo %v: %s", filename, photo, err)
		return
	}

	var b = &bytes.Buffer{}
	// TODO use quality settings from rendition configuration
	err = jpeg.Encode(b, resized, &jpeg.Options{Quality: 95})
	if err != nil {
		log.Printf("Could not resize file %s for photo %v: %s", filename, photo, err)
		return
	}

	backend.Store(record.ID, b.Bytes())
}

func withTransaction(db *sqlx.DB, f func() error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = f()
	if err != nil {
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
