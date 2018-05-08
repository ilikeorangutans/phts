package dbtest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewPhotoRecord(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := CreateCollection(t, dbx)

		repo := db.NewPhotoDB(dbx)

		record := db.PhotoRecord{
			CollectionID: col.ID,
			Filename:     "image.jpg",
			Description:  "it's a photo",
		}
		record, err := repo.Save(record)
		assert.Nil(t, err)
		assert.True(t, record.ID > 0)
	})
}

func TestList(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := CreateCollection(t, dbx)
		col2, _ := CreateCollection(t, dbx)
		repo := db.NewPhotoDB(dbx)

		photo, photoDB := CreatePhoto(t, dbx, col)
		CreatePhoto(t, dbx, col2)

		photo, _ = photoDB.FindByID(col.ID, photo.ID)

		paginator := db.NewPaginator()
		records, err := repo.List(col.ID, paginator)

		assert.Nil(t, err)
		assert.Equal(t, []db.PhotoRecord{photo}, records)
	})
}

func TestListAlbum(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := CreateCollection(t, dbx)
		repo := db.NewPhotoDB(dbx)
		album, albumRepo := CreateAlbum(t, dbx, col)

		photo1, photoRepo := CreatePhoto(t, dbx, col)
		CreatePhoto(t, dbx, col)

		err := albumRepo.AddPhotos(col.ID, album.ID, []int64{photo1.ID})
		photo1, _ = photoRepo.FindByID(col.ID, photo1.ID)

		paginator := db.NewPaginator()
		records, err := repo.ListAlbum(col.ID, album.ID, paginator)

		assert.Nil(t, err)
		assert.Equal(t, []db.PhotoRecord{photo1}, records)
	})
}

func TestListAlbumWithPaginaton(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := CreateCollection(t, dbx)
		repo := db.NewPhotoDB(dbx)
		album, albumRepo := CreateAlbum(t, dbx, col)

		photo1, photoRepo := CreatePhoto(t, dbx, col)
		photo2, _ := CreatePhoto(t, dbx, col)

		err := albumRepo.AddPhotos(col.ID, album.ID, []int64{photo1.ID, photo2.ID})
		assert.Nil(t, err)
		photo1, _ = photoRepo.FindByID(col.ID, photo1.ID)
		photo2, _ = photoRepo.FindByID(col.ID, photo2.ID)

		paginator := db.NewPaginator()
		paginator.PrevID = photo2.ID
		paginator.PrevTimestamp = &photo2.UpdatedAt
		paginator.Direction = db.Asc
		records, err := repo.ListAlbum(col.ID, album.ID, paginator)

		assert.Nil(t, err)
		assert.Equal(t, []db.PhotoRecord{photo1}, records)
	})
}

func TestDeletePhotos(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		col, _ := CreateCollection(t, dbx)

		photo, photoRepo := CreatePhoto(t, dbx, col)

		err := photoRepo.Delete(col.ID, photo.ID)

		assert.Nil(t, err)
	})
}
