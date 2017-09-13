package modeltest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestCreateCollection(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		backend := &storage.FileBackend{BaseDir: "/tmp/backend"}
		repo := model.NewCollectionRepository(dbx, backend)

		col := repo.Create("Test", "test")

		col, err := repo.Save(col)
		assert.Nil(t, err)

		assert.True(t, col.ID > 0)

		repo.Delete(col)
	})
}

func TestAddPhotoToCollectionCreatesRenditions(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		repo := createCollectionRepository(t, dbx)
		col, err := repo.Save(repo.Create("Test", "test"))
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

		// Expecting +1 because we want the original too
		assert.Equal(t, len(renditionConfigs)+1, len(renditions))
	})
}
