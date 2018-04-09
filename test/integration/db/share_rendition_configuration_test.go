package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSetForShare(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		shareSite, _ := CreateShareSite(t, dbx)
		collection, _ := CreateCollection(t, dbx)
		photo, _ := CreatePhoto(t, dbx, collection)
		renditionConfig1, _ := CreateRenditionConfiguration(t, dbx, collection.ID)
		renditionConfig2, _ := CreateRenditionConfiguration(t, dbx, collection.ID)
		share, _ := CreateShare(t, dbx, collection, shareSite, photo, []db.RenditionConfigurationRecord{renditionConfig1})
		repo := db.NewShareRenditionConfigurationDB(dbx)

		_, err := repo.SetForShare(share.ID, []db.ShareRenditionConfigurationRecord{
			{
				RenditionConfigurationID: renditionConfig1.ID,
			},
			{
				RenditionConfigurationID: renditionConfig2.ID,
			},
		})
		assert.Nil(t, err)

		configs, err := repo.FindByShare(share.ID)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(configs))
		for _, config := range configs {
			assert.Equal(t, share.ID, config.ShareID)
		}
	})
}

func TestSetForShareRemoves(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		shareSite, _ := CreateShareSite(t, dbx)
		collection, _ := CreateCollection(t, dbx)
		photo, _ := CreatePhoto(t, dbx, collection)
		renditionConfig1, _ := CreateRenditionConfiguration(t, dbx, collection.ID)
		renditionConfig2, _ := CreateRenditionConfiguration(t, dbx, collection.ID)
		share, _ := CreateShare(t, dbx, collection, shareSite, photo, []db.RenditionConfigurationRecord{renditionConfig1})
		repo := db.NewShareRenditionConfigurationDB(dbx)

		_, err := repo.SetForShare(share.ID, []db.ShareRenditionConfigurationRecord{
			{
				RenditionConfigurationID: renditionConfig1.ID,
			},
		})
		assert.Nil(t, err)

		_, err = repo.SetForShare(share.ID, []db.ShareRenditionConfigurationRecord{
			{
				RenditionConfigurationID: renditionConfig2.ID,
			},
		})
		assert.Nil(t, err)

		configs, err := repo.FindByShare(share.ID)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(configs))
		assert.Equal(t, renditionConfig2.ID, configs[0].RenditionConfigurationID)
	})
}
