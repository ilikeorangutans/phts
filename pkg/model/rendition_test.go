package model

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestInsertRendition(t *testing.T) {
	WithSQLMock(t, func(t *testing.T, ctx context.Context, dbx *sqlx.DB, mock sqlmock.Sqlmock) {
		now := time.Now()
		mock.ExpectQuery("INSERT INTO renditions").
			WithArgs(now, now, 42, true, 1024, 768, "image/jpeg", 17).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(13)).
			RowsWillBeClosed()
		rendition := Rendition{
			Timestamps: db.Timestamps{
				CreatedAt: now,
				UpdatedAt: now,
			},
			Format:                   "image/jpeg",
			Height:                   768,
			Original:                 true,
			PhotoID:                  42,
			RenditionConfigurationID: 17,
			Width:                    1024,
		}

		rendition, err := InsertRendition(ctx, dbx, rendition)

		assert.NoError(t, err)
		assert.Equal(t, int64(13), rendition.ID)
	})
}
