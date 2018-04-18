package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewShareRecord(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		shareSite, _ := CreateShareSite(t, dbx)
		collection, _ := CreateCollection(t, dbx)
		photo, _ := CreatePhoto(t, dbx, collection)
		shareDB := db.NewShareDB(dbx)

		share := db.ShareRecord{
			PhotoID:      photo.ID,
			CollectionID: collection.ID,
			ShareSiteID:  shareSite.ID,
			Slug:         "testing",
		}

		share, err := shareDB.Save(share)
		assert.Nil(t, err)
		assert.True(t, share.IsPersisted())
	})
}

func TestFindByPhoto(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		shareSite, _ := CreateShareSite(t, dbx)
		collection, _ := CreateCollection(t, dbx)
		photo, _ := CreatePhoto(t, dbx, collection)
		shareDB := db.NewShareDB(dbx)
		renditionConfig, _ := CreateRenditionConfiguration(t, dbx, collection.ID)
		shareRendConfDB := db.NewShareRenditionConfigurationDB(dbx)

		share := db.ShareRecord{
			PhotoID:      photo.ID,
			CollectionID: collection.ID,
			ShareSiteID:  shareSite.ID,
			Slug:         "testing",
		}
		share, err := shareDB.Save(share)
		assert.Nil(t, err)
		_, err = shareRendConfDB.SetForShare(share.ID, []db.ShareRenditionConfigurationRecord{{ShareID: share.ID, RenditionConfigurationID: renditionConfig.ID}})
		assert.Nil(t, err)

		result, err := shareRendConfDB.FindByShare(share.ID)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(result))
	})
}

func TestFindByShareSiteAndSlug(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		shareSite, _ := CreateShareSite(t, dbx)
		collection, _ := CreateCollection(t, dbx)
		photo, _ := CreatePhoto(t, dbx, collection)
		shareDB := db.NewShareDB(dbx)
		renditionConfig1, _ := CreateRenditionConfiguration(t, dbx, collection.ID)
		CreateRenditionConfiguration(t, dbx, collection.ID)
		share, err := shareDB.Save(db.ShareRecord{PhotoID: photo.ID, CollectionID: collection.ID, ShareSiteID: shareSite.ID, Slug: "testing"})
		assert.Nil(t, err)
		shareRendConfDB := db.NewShareRenditionConfigurationDB(dbx)
		_, err = shareRendConfDB.SetForShare(share.ID, []db.ShareRenditionConfigurationRecord{{ShareID: share.ID, RenditionConfigurationID: renditionConfig1.ID}})
		assert.Nil(t, err)

		result, err := shareDB.FindByShareSiteAndSlug(shareSite.ID, share.Slug)
		assert.Nil(t, err)

		t.Logf("%s", result)
	})
}
