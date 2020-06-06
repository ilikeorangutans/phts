package model

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func withPhotoRepo(t *testing.T, f func(mock sqlmock.Sqlmock, db *sqlx.DB, repo *PhotoRepo, now time.Time)) {
	now := time.Now()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error mocking connection")
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")

	repo := NewPhotoRepo()

	f(mock, dbx, repo, now)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList(t *testing.T) {
	withPhotoRepo(t, func(mock sqlmock.Sqlmock, dbx *sqlx.DB, repo *PhotoRepo, now time.Time) {
		ctx := context.Background()
		user := User{
			Record: db.Record{ID: 13},
		}

		mock.ExpectQuery("SELECT photos.* FROM photos").
			WillReturnRows(
				sqlmock.
					NewRows(
						[]string{"id", "collection_id", "description", "taken_at", "filename", "rendition_count", "published", "created_at", "updated_at"},
					).
					AddRow(42, 3, "description", nil, "foobar.jpg", 0, true, now, now),
			)

		photos, _, err := repo.List(ctx, dbx, user, database.NewPaginator())

		assert.NoError(t, err)
		assert.Len(t, photos, 1)

		photo := photos[0]
		assert.Equal(t, int64(42), photo.ID)
	})
}

func TestCreatePhoto(t *testing.T) {
	withPhotoRepo(t, func(mock sqlmock.Sqlmock, dbx *sqlx.DB, repo *PhotoRepo, now time.Time) {
		ctx := context.Background()
		photo := Photo{
			Timestamps: db.Timestamps{
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		mock.ExpectQuery("INSERT INTO photos").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(13)).RowsWillBeClosed()

		photo, err := repo.Create(ctx, dbx, photo)

		assert.NoError(t, err)
		assert.Equal(t, int64(13), photo.ID)
	})
}
