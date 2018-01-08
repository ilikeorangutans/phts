package modeltest

import (
	"io/ioutil"
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
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
	return model.NewUserCollectionRepository(dbx, backend, model.User{UserRecord: db.UserRecord{Record: db.Record{ID: 1}}})
}

func createTestCollection(t *testing.T, dbx db.DB) (model.Collection, model.CollectionRepository) {
	repo := createCollectionRepository(t, dbx)
	col := repo.Create("Test Collection", "test-collection")
	col, err := repo.Save(col)
	assert.Nil(t, err)
	return col, repo
}
