package test

import (
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func NewTestDB() (db.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	dbx := sqlx.NewDb(db, "sqlmock")

	return &WrappedDB{
		dbx,
	}, mock
}

type WrappedDB struct {
	*sqlx.DB
}

func (w *WrappedDB) Beginx() (db.TX, error) {
	return w.DB.Beginx()
}
