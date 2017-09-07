package db

import (
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func fixedClock() (time.Time, func() time.Time) {
	now := time.Now()
	return now, func() time.Time { return now }
}

func newTestDB() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	dbx := sqlx.NewDb(db, "sqlmock")

	return dbx, mock
}

func justInsertedRow(id int64) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id"}).AddRow(id)
}

func TestSaveNewRow(t *testing.T) {
	db, mock := newTestDB()

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

func TestUpdate(t *testing.T) {
	db, mock := newTestDB()

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
	mock.ExpectQuery("SELECT count(.+) FROM photos (.+)").WithArgs(
		13,
	).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(17))
	mock.ExpectExec("UPDATE collections").WithArgs(
		"Test", "test", 17, sqlmock.AnyArg(), record.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	record, err := collectionDB.Save(record)

	assert.Nil(t, err)
	assert.True(t, record.IsPersisted())
	assert.Equal(t, int64(13), record.ID)

	assert.Nil(t, mock.ExpectationsWereMet())
}
