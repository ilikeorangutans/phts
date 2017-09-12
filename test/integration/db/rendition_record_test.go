package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewRendition(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, colRepo := createCollection(t, dbx)
		defer colRepo.Delete(col.ID)
		photo, _ := createPhoto(t, dbx, col)

		repo := db.NewRenditionDB(dbx)

		rendition := db.RenditionRecord{
			PhotoID:  photo.ID,
			Original: true,
			Width:    640,
			Height:   480,
			Format:   "image/jpeg",
		}

		rendition, err := repo.Save(rendition)

		assert.Nil(t, err)
		assert.True(t, rendition.ID > 0)
	})
}
