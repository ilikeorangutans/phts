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
