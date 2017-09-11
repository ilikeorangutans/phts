package integration

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func GetDB(t *testing.T) *sqlx.DB {
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

	return dbx
}
