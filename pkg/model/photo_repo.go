package model

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewPhotoRepo() *PhotoRepo {
	return &PhotoRepo{
		clock: time.Now,
		stmt:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type PhotoRepo struct {
	clock func() time.Time
	stmt  sq.StatementBuilderType
}

func (p *PhotoRepo) List(ctx context.Context, db *sqlx.DB, user User, paginator database.Paginator) ([]Photo, database.Paginator, error) {
	stmt := p.stmt.
		Select("photos.*").
		From("photos").
		Join("collections c on (photos.collection_id = c.id)").
		Join("users_collections uc on (uc.collection_id = c.id)").
		Where(sq.Eq{"uc.user_id": user.ID})

	sql, args, err := paginator.Paginate(stmt).ToSql()
	if err != nil {
		return nil, paginator, errors.Wrap(err, "could not build query")
	}

	var photos []Photo
	err = db.SelectContext(ctx, &photos, sql, args...)
	if err != nil {
		return nil, paginator, errors.Wrap(err, "could not select rows")
	}

	return photos, paginator, nil
}

// Create stores a new photo in the database.
func (p *PhotoRepo) Create(ctx context.Context, tx sqlx.ExtContext, photo Photo) (Photo, error) {
	sql, args, err := p.stmt.Insert("photos").
		Columns("updated_at", "created_at", "collection_id", "rendition_count", "description", "filename", "taken_at", "published").
		Values(photo.UpdatedAt, photo.CreatedAt, photo.CollectionID, photo.RenditionCount, photo.Description, photo.Filename, photo.TakenAt, photo.Published).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return photo, errors.Wrap(err, "could not build query")
	}

	rows, err := tx.QueryxContext(ctx, sql, args...)
	defer rows.Close()
	if err != nil {
		return photo, errors.Wrap(err, "could not insert")
	}

	if !rows.Next() {
		return photo, errors.Wrap(err, "no id returned")
	}

	err = rows.Scan(&photo.ID)
	if err != nil {
		return photo, errors.Wrap(err, "could scan id")
	}

	return photo, nil
}

func (p *PhotoRepo) Save(ctx context.Context, tx sqlx.ExtContext, photo Photo) (Photo, error) {
	return photo, nil
}
