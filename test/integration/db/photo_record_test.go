package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewPhotoRecord(t *testing.T) {
	dbx := GetDB(t)
	defer dbx.Close()
	col, colRepo := createCollection(t, dbx)
	defer colRepo.Delete(col.ID)

	repo := db.NewPhotoDB(dbx)

	record := db.PhotoRecord{
		CollectionID: col.ID,
		Filename:     "image.jpg",
		Description:  "it's a photo",
	}
	record, err := repo.Save(record)
	assert.Nil(t, err)
	assert.True(t, record.ID > 0)
}
