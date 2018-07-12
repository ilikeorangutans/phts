package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func fixedClock() (time.Time, func() time.Time) {
	now := time.Now()
	return now, func() time.Time { return now }
}

func justInsertedRow(id int64) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id"}).AddRow(id)
}

func TestSaveNewRow(t *testing.T) {
	db, mock := NewTestDB()

	record := CollectionRecord{
		Sluggable: Sluggable{
			Slug: "test",
		},
		Name: "Test",
	}

	_, clock := fixedClock()
	collectionDB := &collectionSQLDB{
		db:    db,
		clock: clock,
	}
	mock.ExpectQuery("INSERT INTO collections").WithArgs(
		"Test", "test", sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnRows(justInsertedRow(1))

	record, err := collectionDB.Save(record)

	assert.Nil(t, err)
	assert.True(t, record.IsPersisted())
	assert.Equal(t, int64(1), record.ID)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateExistingCollectionRecord(t *testing.T) {
	db, mock := NewTestDB()

	record := CollectionRecord{
		Record: Record{
			ID: 13,
		},
		Sluggable: Sluggable{
			Slug: "test",
		},
		Name: "Test",
	}

	_, clock := fixedClock()
	collectionDB := &collectionSQLDB{
		db:    db,
		clock: clock,
	}
	mock.ExpectExec("UPDATE collections").WithArgs(
		"Test", "test", sqlmock.AnyArg(), record.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	record, err := collectionDB.Save(record)

	assert.Nil(t, err)
	assert.True(t, record.IsPersisted())
	assert.Equal(t, int64(13), record.ID)

	assert.Nil(t, mock.ExpectationsWereMet())
}
