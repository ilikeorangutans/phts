package model

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewCollectionRepo(db *sqlx.DB) (*CollectionRepo, error) {
	return &CollectionRepo{
		db:    db,
		clock: time.Now,
		stmt:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}, nil
}

type CollectionRepo struct {
	db    *sqlx.DB
	clock func() time.Time
	stmt  sq.StatementBuilderType
}

// NewCollection creates a new collection with the given name and slug for the user.
func (c *CollectionRepo) NewCollection(ctx context.Context, name, slug string, owner User) (Collection, error) {
	collection := Collection{
		Timestamps: db.JustCreated(c.clock),
		Sluggable: db.Sluggable{
			Slug: slug,
		},
		Name:       name,
		PhotoCount: 0,
	}

	tx, err := c.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return collection, errors.Wrap(err, "could not begin transaction")
	}
	collection, err = c.create(ctx, tx, collection)
	if err != nil {
		tx.Rollback()
		return collection, errors.Wrap(err, "could not persist collection")
	}
	err = c.createOwner(ctx, tx, owner, collection)
	if err != nil {
		tx.Rollback()
		return collection, errors.Wrap(err, "could not persist collection")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return collection, errors.Wrap(err, "could not commit transaction")
	}

	return collection, nil
}

// createOwner creates an entry in the users_collections table to signify that the given user is an owner of the collection.
func (c *CollectionRepo) createOwner(ctx context.Context, tx sqlx.Ext, user User, collection Collection) error {
	result, err := c.stmt.
		Insert("users_collections").
		Columns("user_id", "collection_id", "created_at", "updated_at").
		Values(user.ID, collection.ID, collection.CreatedAt, collection.UpdatedAt).
		RunWith(c.db).
		ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "could not persist ownership record")
	}

	if rows, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "could not get number of affected rows")
	} else if rows != 1 {
		return errors.Wrap(err, "did not update one row")
	}

	return nil
}

func (c *CollectionRepo) create(ctx context.Context, tx sqlx.Ext, collection Collection) (Collection, error) {
	err := c.stmt.
		Insert("collections").
		Columns("name", "slug", "created_at", "updated_at", "photo_count").
		Values(collection.Name, collection.Slug, collection.CreatedAt, collection.UpdatedAt, 0).
		Suffix("returning id").
		RunWith(c.db).
		QueryRowContext(ctx).
		Scan(&collection.ID)
	if err != nil {
		return collection, errors.Wrap(err, "could not persist collection")
	}

	return collection, nil
}
