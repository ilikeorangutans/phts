package modeltest

import (
	"io/ioutil"
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateCollection(t *testing.T) {
	dbx := GetDB(t)
	defer dbx.Close()

	backend := &storage.FileBackend{BaseDir: "/tmp/backend"}
	repo := model.NewCollectionRepository(dbx, backend)

	col := repo.Create("Test", "test")

	col, err := repo.Save(col)
	assert.Nil(t, err)

	assert.True(t, col.ID > 0)

	repo.Delete(col)
}

func get1x1JPEG(t *testing.T) []byte {
	b, err := ioutil.ReadFile("../files/1x1.jpg")
	assert.Nil(t, err)
	return b
}

func createCollectionRepository(t *testing.T, dbx *sqlx.DB) model.CollectionRepository {
	name, err := ioutil.TempDir("", "file-backend")
	assert.Nil(t, err)
	backend := &storage.FileBackend{BaseDir: name}
	return model.NewCollectionRepository(dbx, backend)
}

func TestAddPhotoToCollectionCreatesRenditions(t *testing.T) {
	dbx := GetDB(t)
	defer dbx.Close()
	repo := createCollectionRepository(t, dbx)
	col, err := repo.Save(repo.Create("Test", "test"))
	assert.Nil(t, err)
	defer repo.Delete(col)

	err = repo.AddPhoto(col, "image.jpg", get1x1JPEG(t))
	assert.Nil(t, err)

	photoDB := db.NewPhotoDB(dbx)
	photos, err := photoDB.List(col.ID, db.NewPaginator())
	assert.Nil(t, err)
	photoRecord := photos[0]

	assert.Equal(t, col.ID, photoRecord.CollectionID)

	renditionDB := db.NewRenditionDB(dbx)
	renditions, err := renditionDB.FindAllForPhoto(photoRecord.ID)
	assert.Nil(t, err)

	renditionConfigDB := db.NewRenditionConfigurationDB(dbx)
	renditionConfigs, err := renditionConfigDB.FindForCollection(col.ID)
	assert.Nil(t, err)

	// Expecting +1 because we want the original too
	assert.Equal(t, len(renditionConfigs)+1, len(renditions))
}
