package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewCollectionRecord(t *testing.T) {
	dbx := integration.GetDB(t)
	defer dbx.Close()
	repo := db.NewCollectionDB(dbx)

	record := db.CollectionRecord{
		Sluggable: db.Sluggable{Slug: "test"},
		Name:      "Test",
	}

	result, err := repo.Save(record)
	assert.Nil(t, err)
	assert.True(t, result.ID > 0)

	err = repo.Delete(result.ID)
	assert.Nil(t, err)
}

func TestUpdateCollectionRecord(t *testing.T) {
	dbx := integration.GetDB(t)
	defer dbx.Close()
	repo := db.NewCollectionDB(dbx)

	record := db.CollectionRecord{
		Sluggable: db.Sluggable{Slug: "test"},
		Name:      "Test",
	}
	record, err := repo.Save(record)
	assert.Nil(t, err)
	defer repo.Delete(record.ID)

	record.Name = "Test Updated"
	record.Slug = "test-updated"

	record, err = repo.Save(record)

	assert.Nil(t, err)
}
