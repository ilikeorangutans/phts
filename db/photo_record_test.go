package db

import (
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/assert"
)

func TestListEmpty(t *testing.T) {
	db, mock := newTestDB()
	_, clock := fixedClock()

	photoDB := &photoSQLDB{
		db:    db,
		clock: clock,
	}

	mock.ExpectQuery("SELECT (.+) FROM photos WHERE collection_id (.+)").WithArgs(13, 10).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	photos, err := photoDB.List(13, NewPaginator())
	assert.Nil(t, err)

	assert.Equal(t, []PhotoRecord{}, photos)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestListReturnsPhotos(t *testing.T) {
	db, mock := newTestDB()
	_, clock := fixedClock()

	photoDB := &photoSQLDB{
		db:    db,
		clock: clock,
	}

	mock.ExpectQuery(
		"SELECT (.+) FROM photos WHERE collection_id (.+)",
	).WithArgs(13, 10).WillReturnRows(
		sqlmock.NewRows([]string{"id", "collection_id", "filename"}).AddRow(11, 13, "image.jpg"),
	)

	photos, err := photoDB.List(13, NewPaginator())
	assert.Nil(t, err)

	assert.Equal(t, 1, len(photos))
	photo := photos[0]
	assert.Equal(t, int64(13), photo.CollectionID)
	assert.Equal(t, "image.jpg", photo.Filename)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSaveNewRecordWithoutCollectionIDFails(t *testing.T) {
	db, _ := newTestDB()
	_, clock := fixedClock()

	photoDB := &photoSQLDB{
		db:    db,
		clock: clock,
	}

	record := PhotoRecord{}

	_, err := photoDB.Save(record)
	assert.NotNil(t, err)
}

func TestSaveNewRecord(t *testing.T) {
	db, mock := newTestDB()
	_, clock := fixedClock()

	photoDB := &photoSQLDB{
		db:    db,
		clock: clock,
	}

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
