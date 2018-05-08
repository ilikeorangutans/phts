package dbtest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/stretchr/testify/assert"
)

func CreateUser(t *testing.T, dbx db.DB) (db.UserRecord, db.UserDB) {
	userDB := db.NewUserDB(dbx)

	user := db.UserRecord{
		Email: "test@test.com",
	}

	user, err := userDB.Save(user)
	assert.Nil(t, err)
	return user, userDB
}

func CreateCollection(t *testing.T, dbx db.DB) (db.CollectionRecord, db.CollectionDB) {
	colRepo := db.NewCollectionDB(dbx)
	col := db.CollectionRecord{
		Sluggable: db.Sluggable{Slug: fmt.Sprintf("test-%d", rand.Int63())},
		Name:      "Test",
	}
	col, err := colRepo.Save(col)
	assert.Nil(t, err)
	return col, colRepo
}

func CreatePhoto(t *testing.T, dbx db.DB, collection db.CollectionRecord) (db.PhotoRecord, db.PhotoDB) {
	repo := db.NewPhotoDB(dbx)
	record := db.PhotoRecord{
		CollectionID: collection.ID,
		Filename:     "image.jpg",
		Description:  "it's a photo",
	}
	record, err := repo.Save(record)
	assert.Nil(t, err)
	return record, repo
}

func CreateAlbum(t *testing.T, dbx db.DB, collection db.CollectionRecord) (db.AlbumRecord, db.AlbumDB) {
	repo := db.NewAlbumDB(dbx)
	record := db.AlbumRecord{
		Name:         "test",
		Slug:         "test",
		CollectionID: collection.ID,
		PhotoCount:   0,
	}
	record, err := repo.Save(record)
	assert.Nil(t, err)
	return record, repo
}

func CreateRenditionConfiguration(t *testing.T, dbx db.DB, collectionID int64) (db.RenditionConfigurationRecord, db.RenditionConfigurationDB) {
	repo := db.NewRenditionConfigurationDB(dbx)
	record := db.RenditionConfigurationRecord{
		Width:        320,
		Height:       240,
		Name:         fmt.Sprintf("test-%d", rand.Int63()),
		Quality:      80,
		CollectionID: &collectionID,
	}

	record, err := repo.Save(record)
	assert.Nil(t, err)

	return record, repo
}

func CreateRenditions(t *testing.T, dbx db.DB, photo db.PhotoRecord, configs []db.RenditionConfigurationRecord) ([]db.RenditionRecord, db.RenditionDB) {
	renditionDB := db.NewRenditionDB(dbx)

	var results []db.RenditionRecord
	for _, config := range configs {

		result, err := renditionDB.Save(db.RenditionRecord{
			PhotoID:                  photo.ID,
			RenditionConfigurationID: config.ID,
			Width:    uint(config.Width),
			Height:   uint(config.Height),
			Format:   "image/jpeg",
			Original: config.Original,
		})
		results = append(results, result)
		assert.Nil(t, err)
	}

	return results, renditionDB
}

func CreateShareSite(t *testing.T, dbx db.DB) (db.ShareSiteRecord, db.ShareSiteDB) {
	repo := db.NewShareSiteDB(dbx)
	record := db.ShareSiteRecord{
		Domain: fmt.Sprintf("%s.phts.org", time.Now().Format("20060102150405.000")),
	}

	record, err := repo.Save(record)
	assert.Nil(t, err)

	return record, repo
}

func CreateShare(t *testing.T, dbx db.DB, collection db.CollectionRecord, shareSite db.ShareSiteRecord, photo db.PhotoRecord) (db.ShareRecord, db.ShareDB) {
	slug, _ := model.SlugFromString(time.Now().Format(time.RFC822Z))
	repo := db.NewShareDB(dbx)
	record := db.ShareRecord{
		PhotoID:      photo.ID,
		CollectionID: collection.ID,
		ShareSiteID:  shareSite.ID,
		Slug:         slug,
	}

	record, err := repo.Save(record)
	assert.Nil(t, err)

	return record, repo
}
