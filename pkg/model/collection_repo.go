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

// FindByIDAndUser finds the collection with the given id for the specified user.
func (c *CollectionRepo) FindByIDAndUser(ctx context.Context, db sqlx.QueryerContext, id int64, user User) (Collection, error) {
	var collection Collection
	sql, args, err := c.stmt.Select("collections.*").
		From("collections").
		Join("users_collections on (users_collections.collection_id = collections.id)").
		Where(sq.Eq{
			"users_collections.user_id": user.ID,
			"collections.id":            id,
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return collection, errors.Wrap(err, "could not create query")
	}

	err = db.QueryRowxContext(ctx, sql, args...).StructScan(&collection)
	if err != nil {
		return collection, errors.Wrap(err, "could not execute query")
	}

	return collection, nil
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

// Update updates the given collection.
func (c *CollectionRepo) Update(ctx context.Context, tx sqlx.Ext, collection Collection) (Collection, error) {
	var photoCount = 0
	c.stmt.Select("count(*)").
		From("photos").
		Where(sq.Eq{"collection_id": collection.ID}).
		RunWith(tx).
		ScanContext(ctx, &photoCount)

	collection.JustUpdated(c.clock)
	result, err := c.stmt.
		Update("collections").
		Set("name", collection.Name).
		Set("slug", collection.Slug).
		Set("updated_at", collection.UpdatedAt).
		Set("photo_count", photoCount).
		RunWith(c.db).
		ExecContext(ctx)
	if err != nil {
		return collection, errors.Wrap(err, "could not update collection")
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return collection, errors.Wrap(err, "could not get number of affected rows")
	} else if rowsAffected != 1 {
		return collection, errors.Wrap(err, "row not updated")
	}

	return collection, nil
}

func (c *CollectionRepo) AddPhotos(ctx context.Context, dbx *sqlx.DB, collection Collection, photoUploads ...PhotoUpload) (Collection, []Photo, error) {
	tx, err := dbx.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return collection, nil, errors.Wrap(err, "could not start transaction")
	}
	photoRepo := NewPhotoRepo()

	var photos []Photo
	for i, upload := range photoUploads {
		photo, err := photoRepo.AddPhoto(ctx, tx, collection, upload)
		if err != nil {
			tx.Rollback()
			return collection, nil, errors.Wrapf(err, "could not add photo %d/%d", i+1, len(photoUploads))
		}

		photos = append(photos, photo)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return collection, nil, errors.Wrap(err, "could not commit transaction")
	}
	return collection, photos, nil
}
