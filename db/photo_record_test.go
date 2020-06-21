package db

import (
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestListReturnsPhotos(t *testing.T) {
	db, mock := NewTestDB()
	_, clock := fixedClock()

	photoDB := NewPhotoDBWithClock(db, clock)

	mock.ExpectQuery(
		"SELECT (.+) FROM photos WHERE collection_id (.+)",
	).WithArgs(13).WillReturnRows(
		sqlmock.NewRows([]string{"id", "collection_id", "filename"}).AddRow(11, 13, "image.jpg"),
	)

	photos, err := photoDB.List(13, database.NewPaginator())
	assert.Nil(t, err)

	assert.Equal(t, 1, len(photos))
	photo := photos[0]
	assert.Equal(t, int64(13), photo.CollectionID)
	assert.Equal(t, "image.jpg", photo.Filename)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSaveNewRecordWithoutCollectionIDFails(t *testing.T) {
	db, _ := NewTestDB()
	_, clock := fixedClock()

	photoDB := NewPhotoDBWithClock(db, clock)

	record := PhotoRecord{}

	_, err := photoDB.Save(record)
	assert.NotNil(t, err)
}

func TestSaveNewRecord(t *testing.T) {
	db, mock := NewTestDB()
	_, clock := fixedClock()

	photoDB := NewPhotoDBWithClock(db, clock)

	now := clock()
	record := PhotoRecord{
		CollectionID: 13,
		Filename:     "image.jpg",
		Description:  "description",
	}

	mock.ExpectQuery(
		"INSERT INTO photos",
	).WithArgs(
		13, "image.jpg", sqlmock.AnyArg(), now.UTC(), now.UTC(),
	).WillReturnRows(
		justInsertedRow(17),
	)

	photo, err := photoDB.Save(record)
	assert.Nil(t, err)
	assert.NotNil(t, photo)

	assert.Equal(t, 0, photo.RenditionCount)
	assert.Equal(t, int64(17), photo.ID)

	if err = mock.ExpectationsWereMet(); err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestUpdateExistingRecord(t *testing.T) {
	// TODO sometimes this test is flaky. Not sure if it's sqlmock or if I'm setting it up wrong.
	db, mock := NewTestDB()
	_, clock := fixedClock()

	photoDB := NewPhotoDBWithClock(db, clock)

	now := clock()
	earlier := now.Add(time.Hour * -1)
	record := PhotoRecord{
		Record: Record{
			ID: 17,
		},
		Timestamps: Timestamps{
			CreatedAt: earlier,
			UpdatedAt: earlier,
		},
		CollectionID: 13,
		Filename:     "image.jpg",
		Description:  "description",
	}

	mock.ExpectExec(
		"UPDATE photos",
	).WithArgs(
		"image.jpg", now.UTC(), 17, 13,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	photo, err := photoDB.Save(record)
	assert.Nil(t, err)

	assert.Equal(t, 0, photo.RenditionCount)
	assert.Equal(t, now, photo.UpdatedAt)

	if err = mock.ExpectationsWereMet(); err != nil {
		assert.Fail(t, err.Error())
	}
}
