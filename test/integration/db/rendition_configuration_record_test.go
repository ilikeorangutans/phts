package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
	"github.com/stretchr/testify/assert"
)

func TestFindRenditionConfigByIDAndCollection(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		collectionDB := db.NewCollectionDB(dbx)
		collectionRecord := db.CollectionRecord{
			Sluggable: db.Sluggable{Slug: "test"},
			Name:      "Test",
		}
		collectionRecord, err := collectionDB.Save(collectionRecord)

		configDB := db.NewRenditionConfigurationDB(dbx)

		record := db.RenditionConfigurationRecord{
			Width:        640,
			Height:       480,
			Name:         "VGA",
			Quality:      85,
			CollectionID: &collectionRecord.ID,
		}

		expected, err := configDB.Save(record)
		assert.Nil(t, err)

		result, err := configDB.FindByID(collectionRecord.ID, expected.ID)
		assert.Nil(t, err)
		assert.Equal(t, expected.ID, result.ID)
	})
}

func TestFindRenditionConfigByIDWithoutCollection(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		collectionDB := db.NewCollectionDB(dbx)
		collectionRecord := db.CollectionRecord{
			Sluggable: db.Sluggable{Slug: "test"},
			Name:      "Test",
		}
		collectionRecord, err := collectionDB.Save(collectionRecord)

		configDB := db.NewRenditionConfigurationDB(dbx)

		record := db.RenditionConfigurationRecord{
			Width:   640,
			Height:  480,
			Name:    "VGA",
			Quality: 85,
		}

		expected, err := configDB.Save(record)
		assert.Nil(t, err)

		result, err := configDB.FindByID(collectionRecord.ID, expected.ID)
		assert.Nil(t, err)
		assert.Equal(t, expected.ID, result.ID)
	})
}
func TestFindAllRenditionConfigs(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		configDB := db.NewRenditionConfigurationDB(dbx)

		configs, err := configDB.FindForCollection(0)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(configs))

		for _, c := range configs {
			t.Logf("- %s", c.Name)
		}
	})
}

func TestSaveNewRenditionConfig(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		configDB := db.NewRenditionConfigurationDB(dbx)

		configs, err := configDB.FindForCollection(0)
		assert.Nil(t, err)
		initialCount := len(configs)

		record := db.RenditionConfigurationRecord{
			Width:   640,
			Height:  480,
			Name:    "VGA",
			Quality: 85,
		}

		result, err := configDB.Save(record)
		assert.Nil(t, err)
		assert.True(t, result.ID > 0)

		configs, err = configDB.FindForCollection(1)
		assert.Nil(t, err)
		assert.Equal(t, initialCount+1, len(configs))

		err = configDB.Delete(result.ID)
		assert.Nil(t, err)
	})
}
