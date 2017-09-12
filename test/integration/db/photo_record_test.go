package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewPhotoRecord(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := createCollection(t, dbx)

		repo := db.NewPhotoDB(dbx)

		record := db.PhotoRecord{
			CollectionID: col.ID,
			Filename:     "image.jpg",
			Description:  "it's a photo",
		}
		record, err := repo.Save(record)
		assert.Nil(t, err)
		assert.True(t, record.ID > 0)
	})
}
