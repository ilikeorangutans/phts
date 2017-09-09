package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
	"github.com/stretchr/testify/assert"
)

func TestFindAllRenditionConfigs(t *testing.T) {
	dbx := GetDB(t)
	defer dbx.Close()
	configDB := db.NewRenditionConfigurationDB(dbx)

	configs, err := configDB.FindForCollection(0)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(configs))

	for _, c := range configs {
		t.Logf("- %s", c.Name)
	}
}

func TestSaveNewRenditionConfig(t *testing.T) {
	dbx := GetDB(t)
	defer dbx.Close()
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

	configs, err = configDB.FindForCollection(0)
	assert.Nil(t, err)
	assert.Equal(t, initialCount+1, len(configs))

	err = configDB.Delete(result.ID)
	assert.Nil(t, err)

}
