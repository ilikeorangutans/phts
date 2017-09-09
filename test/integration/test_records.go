package integration

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func createCollection(t *testing.T, dbx *sqlx.DB) (db.CollectionRecord, db.CollectionDB) {
	colRepo := db.NewCollectionDB(dbx)
	col := db.CollectionRecord{
		Sluggable: db.Sluggable{Slug: "test"},
		Name:      "Test",
	}
	col, err := colRepo.Save(col)
	assert.Nil(t, err)
	return col, colRepo
}

func createPhoto(t *testing.T, dbx *sqlx.DB, collection db.CollectionRecord) (db.PhotoRecord, db.PhotoDB) {
	repo := db.NewPhotoDB(dbx)
	record := db.PhotoRecord{
		CollectionID: collection.ID,
		Filename:     "image.jpg",
		Description:  "it's a photo",
	}
	record, err := repo.Save(record)
	assert.Nil(t, err)
	return record, repo
}

func createRenditionConfiguration(t *testing.T, dbx *sqlx.DB, collectionID int64) (db.RenditionConfigurationRecord, db.RenditionConfigurationDB) {
	repo := db.NewRenditionConfigurationDB(dbx)
	record := db.RenditionConfigurationRecord{
		Width:        320,
		Height:       240,
		Name:         "test",
		Quality:      80,
		CollectionID: &collectionID,
	}

	record, err := repo.Save(record)
	assert.Nil(t, err)

	return record, repo
}
