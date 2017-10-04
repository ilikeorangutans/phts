package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewRendition(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := createCollection(t, dbx)
		renditionConfigDB := db.NewRenditionConfigurationDB(dbx)
		config, _ := renditionConfigDB.FindByName(0, "original")
		photo, _ := createPhoto(t, dbx, col)

		repo := db.NewRenditionDB(dbx)

		rendition := db.RenditionRecord{
			PhotoID:  photo.ID,
			Original: true,
			Width:    640,
			Height:   480,
			Format:   "image/jpeg",
			RenditionConfigurationID: config.ID,
		}

		rendition, err := repo.Save(rendition)

		assert.Nil(t, err)
		assert.True(t, rendition.ID > 0)
	})
}

func TestFindByPhotoAndConfigs(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := createCollection(t, dbx)
		photo, _ := createPhoto(t, dbx, col)

		repo := db.NewRenditionDB(dbx)

		_, err := repo.Save(db.RenditionRecord{
			PhotoID:  photo.ID,
			Original: false,
			Width:    640,
			Height:   480,
			Format:   "image/jpeg",
			RenditionConfigurationID: 1,
		})
		assert.Nil(t, err)

		result, err := repo.FindByPhotoAndConfigs(col.ID, photo.ID, []int64{1})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})
}
