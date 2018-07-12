package modeltest

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
	testdb "github.com/ilikeorangutans/phts/test/integration/db"
	"github.com/stretchr/testify/assert"
)

func get1x1JPEG(t *testing.T) []byte {
	b, err := ioutil.ReadFile("../files/1x1.jpg")
	assert.Nil(t, err)
	return b
}

func getSmallJPEGWithExif(t *testing.T) []byte {
	b, err := ioutil.ReadFile("../files/100x75-with-exif.jpg")
	assert.Nil(t, err)
	return b
}

func getStorage(t *testing.T) storage.Backend {
	name, err := ioutil.TempDir("", "file-backend")
	assert.Nil(t, err)
	return &storage.FileBackend{BaseDir: name}
}

func createCollectionRepository(t *testing.T, dbx db.DB) model.CollectionRepository {
	backend := getStorage(t)
	user, _ := testdb.CreateUser(t, dbx)
	return model.NewUserCollectionRepository(dbx, backend, user)
}

func createTestCollection(t *testing.T, dbx db.DB) (model.Collection, model.CollectionRepository) {
	repo := createCollectionRepository(t, dbx)
	col := repo.Create("Test Collection", "test-collection")
	col, err := repo.Save(col)
	assert.Nil(t, err)
	return col, repo
}

func CreatePhoto(t *testing.T, dbx db.DB, col model.Collection) (model.Photo, model.PhotoRepository) {
	backend := getStorage(t)
	repo := model.NewPhotoRepository(dbx, backend)

	photo, err := repo.Create(col, fmt.Sprintf("img-%s.jpg", time.Now().Format("20060102150405.000")), get1x1JPEG(t))
	assert.Nil(t, err)

	return photo, repo
}

func CreateShareSite(t *testing.T, dbx db.DB) (model.ShareSite, model.ShareSiteRepository) {
	record, _ := testdb.CreateShareSite(t, dbx)

	repo := model.NewShareSiteRepository(dbx)

	shareSite, err := repo.Save(model.ShareSite{record})
	assert.Nil(t, err)

	return shareSite, repo
}

func CreateRenditionConfigurations(t *testing.T, dbx db.DB, col model.Collection) (model.RenditionConfigurations, model.RenditionConfigurationRepository) {
	create := []struct {
		name     string
		original bool
	}{
		{"small", false},
		{"medium", false},
		{"original", true},
	}

	repo := model.NewRenditionConfigurationRepository(dbx)
	var result model.RenditionConfigurations
	for _, data := range create {
		config := model.RenditionConfiguration{
			db.RenditionConfigurationRecord{
				Name: data.name,
			},
		}

		config, err := repo.Save(col, config)
		assert.Nil(t, err)

		result = append(result, config)
	}

	return result, repo
}
