package model

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
)

func TestFindMissingrenditionConfigurations(t *testing.T) {
	WithSQLMock(t, func(t *testing.T, ctx context.Context, dbx *sqlx.DB, mock sqlmock.Sqlmock) {
		FindMissingRenditionConfigurations(ctx, dbx, Photo{Record: db.Record{ID: 13}})
		// TODO finish me
	})
}
