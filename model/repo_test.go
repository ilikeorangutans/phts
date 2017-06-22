package model

import (
	"testing"

	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFoox(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	dbx := sqlx.NewDb(db, "sqlmock")
	repo := &CollectionSQLRepository{
		db: dbx,
	}

	mock.ExpectQuery("SELECT * FROM collections WHERE id=$1").WithArgs(1)

	repo.FindByID(1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
