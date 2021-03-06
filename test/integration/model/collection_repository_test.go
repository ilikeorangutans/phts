package modeltest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/test/integration"
	dbtest "github.com/ilikeorangutans/phts/test/integration/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateCollection(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		backend := &storage.FileBackend{BaseDir: "/tmp/backend"}
		user, _ := dbtest.CreateUser(t, dbx)
		repo := model.NewUserCollectionRepository(dbx, backend, user)

		col := repo.NewInstance("Test", "test")

		err := repo.Save(col)
		assert.Nil(t, err)

		assert.True(t, col.ID > 0)

		repo.Delete(col)
	})
}

func TestAddPhotoToCollectionCreatesRenditions(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		repo := createCollectionRepository(t, dbx)
		col := repo.NewInstance("Test", "test")
		err := repo.Save(col)
		assert.Nil(t, err)
		defer repo.Delete(col)

		_, err = repo.AddPhoto(col, "image.jpg", get1x1JPEG(t))
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

		assert.Equal(t, len(renditionConfigs), len(renditions))
	})
}
