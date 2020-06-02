package model

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func WithSQLMock(t *testing.T, f func(t *testing.T, ctx context.Context, dbx *sqlx.DB, mock sqlmock.Sqlmock)) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error mocking connection")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbx := sqlx.NewDb(db, "postgres")

	f(t, ctx, dbx, mock)

	assert.NoError(t, mock.ExpectationsWereMet())
}
