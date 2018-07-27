package integration

import (
	"log"
	"testing"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	if err := dbx.DB.Ping(); err != nil {
		log.Printf("Error connecting: %v", err.Error())
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
	tx, err := dbx.Begin()
	if err != nil {
		t.Log("Error while starting transaction for migration: %s", err.Error())
		t.Fail()
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
	} else if err != nil {
		t.Log("Error while migrating database: %s", err.Error())
		t.Fail()
	}
	err = tx.Commit()
	if err != nil {
		t.Log("Error while starting transaction for migration: %s", err.Error())
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
