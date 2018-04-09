package modeltest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSaveShare(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, colRepo := createTestCollection(t, dbx)
		shareSite, _ := CreateShareSite(t, dbx)
		configs, _ := CreateRenditionConfigurations(t, dbx, col)
		photo, _ := CreatePhoto(t, dbx, col)
		repo := model.NewShareRepository(dbx, colRepo)

		share, errors := shareSite.Builder().
			FromCollection(col).
			AddPhoto(photo).
			WithSlug("I am a slug!").
			AllowRenditions(configs).
			Build()

		assert.True(t, len(errors) == 0)

		share, err := repo.Publish(share)
		assert.Nil(t, err)

		shares, err := repo.FindByPhoto(photo, db.Paginator{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(shares))
		assert.Equal(t, "i-am-a-slug-", shares[0].Slug)
		assert.Equal(t, 3, len(shares[0].RenditionConfigurations))
	})
}
