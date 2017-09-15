package modeltest

import (
	"testing"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestPhotoRepositoryCreate(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := createTestCollection(t, dbx)
		backend := getStorage(t)
		repo := model.NewPhotoRepository(dbx, backend)

		photo, err := repo.Create(col, "image.jpg", get1x1JPEG(t))
		assert.Nil(t, err)
		assert.Equal(t, col.ID, photo.CollectionID)

		renditionRepo := db.NewRenditionDB(dbx)
		renditions, err := renditionRepo.FindAllForPhoto(photo.ID)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(renditions))

		renditionConfigs := db.NewRenditionConfigurationDB(dbx)
		for _, r := range renditions {
			data, err := backend.Get(r.ID)
			assert.Nil(t, err)
			config, err := renditionConfigs.FindByID(col.ID, r.RenditionConfigurationID)
			assert.Nil(t, err)

			t.Logf("rendition %d, original: %t, size: %d, config: %d %s", r.ID, r.Original, len(data), config.ID, config.Name)
			assert.True(t, r.Original == (config.Name == "original"))
		}
	})
}

func TestPhotoRepositoryCreateCheckExif(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := createTestCollection(t, dbx)
		backend := getStorage(t)
		repo := model.NewPhotoRepository(dbx, backend)

		photo, err := repo.Create(col, "image.jpg", getSmallJPEGWithExif(t))
		assert.Nil(t, err)

		assert.Equal(t, time.Date(2015, time.August, 1, 19, 50, 0, 0, time.UTC), *photo.TakenAt)
	})
}
