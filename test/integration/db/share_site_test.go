package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewShareSiteRecord(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		repo := db.NewShareSiteDB(dbx)

		record := db.ShareSiteRecord{
			Domain: "photos.phts.org",
		}

		result, err := repo.Save(record)

		assert.Nil(t, err)
		assert.NotNil(t, result.ID)

	})
}

func TestUpdateShareSiteRecord(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {

		repo := db.NewShareSiteDB(dbx)

		record := db.ShareSiteRecord{
			Domain: "photos.phts.org",
		}

		record, err := repo.Save(record)

		assert.Nil(t, err)
		assert.True(t, record.IsPersisted())

		record.Domain = "photos2.phts.org"

		record, err = repo.Save(record)
		assert.Nil(t, err)

		assert.Equal(t, record.Domain, "photos2.phts.org")

	})
}
