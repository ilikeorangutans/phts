package model

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestFindByIDAndUserNoResults(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error mocking connection")
	}
	defer db.Close()
	dbx := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery("").RowsWillBeClosed().WillReturnRows(sqlmock.NewRows([]string{"id"}))

	repo, _ := NewCollectionRepo(dbx)
	ctx := context.Background()

	_, err = repo.FindByIDAndUser(ctx, dbx, 13, User{})

	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func TestFindByIDAndUser(t *testing.T) {
	dbmock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error mocking connection")
	}
	defer dbmock.Close()
	dbx := sqlx.NewDb(dbmock, "postgres")

	mock.ExpectQuery("SELECT collections.* FROM collections JOIN users_collections").
		WithArgs(int64(13), int64(17)).
		RowsWillBeClosed().
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at", "slug", "name", "photo_count"}).
				AddRow(int64(13), time.Now(), time.Now(), "slug", "test", 42),
		)

	repo, _ := NewCollectionRepo(dbx)
	ctx := context.Background()

	collection, err := repo.FindByIDAndUser(ctx, dbx, int64(13), User{Record: db.Record{ID: 17}})

	assert.NoError(t, err)
	assert.Equal(t, int64(13), collection.ID)
}
