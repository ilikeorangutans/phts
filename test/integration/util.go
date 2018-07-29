package integration

import (
	"log"
	"testing"

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

	if err := db.ApplyMigrations(dbx.DB); err != nil {
		t.Logf("Failed to apply migrations: %s", err)
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
