package integration

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

var databaseInstance *sqlx.DB

func createDB(t *testing.T) *sqlx.DB {
	if databaseInstance != nil {
		return databaseInstance
	}

	dbx, err := sqlx.Open("postgres", "user=phts_test password=phts dbname=phts_test sslmode=disable")
	if err != nil {
		t.Log("Error while connecting to postgres: %s", err.Error())
		t.Fail()
	}

	driver, err := postgres.WithInstance(dbx.DB, &postgres.Config{})
	if err != nil {
		t.Log("Error while getting driver: %s", err.Error())
		t.Fail()
	}
	m, err := migrate.NewWithDatabaseInstance("file://../../../db/migrate", "postgres", driver)
	if err != nil {
		t.Log("Error while creating migration: %s", err.Error())
		t.Fail()
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
	} else if err != nil {
		t.Log("Error while migrating database: %s", err.Error())
		t.Fail()
	}

	databaseInstance = dbx
	return dbx
}

func RunTestInDB(t *testing.T, f func(dbx db.DB)) {
	dbx := createDB(t)
	tx, err := dbx.Beginx()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
		return
	}

	wrappedDB := &TXAsDBWrapper{
		tx: tx,
	}
	f(wrappedDB)

	if err := tx.Rollback(); err != nil {
		// Not necessarily a failure
		t.Log(err.Error())
	}
}
