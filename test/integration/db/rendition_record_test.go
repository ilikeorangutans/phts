package dbtest

import (
	"database/sql"
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewRendition(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := CreateCollection(t, dbx)
		renditionConfigDB := db.NewRenditionConfigurationDB(dbx)
		config, _ := renditionConfigDB.FindByName(0, "original")
		photo, _ := CreatePhoto(t, dbx, col)

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
		col, _ := CreateCollection(t, dbx)
		photo, _ := CreatePhoto(t, dbx, col)

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

func TestFindByShareAndIDReturnsRenditionIfAllowedByConfig(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := CreateCollection(t, dbx)
		photo, _ := CreatePhoto(t, dbx, col)
		renditionConfig, _ := CreateRenditionConfiguration(t, dbx, col.ID)
		shareSite, _ := CreateShareSite(t, dbx)
		share, _ := CreateShare(t, dbx, col, shareSite, photo)
		shareRenditionConfigDB := db.NewShareRenditionConfigurationDB(dbx)
		shareRenditionConfigDB.SetForShare(share.ID, []db.ShareRenditionConfigurationRecord{{ShareID: share.ID, RenditionConfigurationID: renditionConfig.ID}})
		renditions, repo := CreateRenditions(t, dbx, photo, []db.RenditionConfigurationRecord{renditionConfig})

		result, err := repo.FindByShareAndID(share.ID, renditions[0].ID)
		assert.Equal(t, renditions[0].ID, result.ID)
		assert.Nil(t, err)
	})
}

func TestFindByShareAndIDDoesNotReturnIfConfigNotAllowed(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := CreateCollection(t, dbx)
		photo, _ := CreatePhoto(t, dbx, col)
		renditionConfig, _ := CreateRenditionConfiguration(t, dbx, col.ID)
		shareSite, _ := CreateShareSite(t, dbx)
		share, _ := CreateShare(t, dbx, col, shareSite, photo)
		renditions, repo := CreateRenditions(t, dbx, photo, []db.RenditionConfigurationRecord{renditionConfig})

		_, err := repo.FindByShareAndID(share.ID, renditions[0].ID)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}
