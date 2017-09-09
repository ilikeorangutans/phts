package dbtest

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/stretchr/testify/assert"
)

func GetDB(t *testing.T) *sqlx.DB {
	dbx, err := sqlx.Open("postgres", "user=dev dbname=phts_test sslmode=disable")
	assert.Nil(t, err)

	driver, err := postgres.WithInstance(dbx.DB, &postgres.Config{})
	assert.Nil(t, err)
	m, err := migrate.NewWithDatabaseInstance("file://../../../db/migrate", "postgres", driver)
	assert.Nil(t, err)
	err = m.Up()
	assert.True(t, err == nil || err == migrate.ErrNoChange)

	return dbx
}
