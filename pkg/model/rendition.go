package model

import (
	"context"
	"time"

	"github.com/ilikeorangutans/phts/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Rendition struct {
	db.Record
	db.Timestamps

	Format                   string `db:"format" json:"format"`
	Height                   uint   `db:"height" json:"height"`
	Original                 bool   `db:"original" json:"original"`
	PhotoID                  int64  `db:"photo_id" json:"photoID"`
	RenditionConfigurationID int64  `db:"rendition_configuration_id" json:"renditionConfigurationID"`
	Width                    uint   `db:"width" json:"width"`
}

// CreateRendition creates a new rendition in the database.
func InsertRendition(ctx context.Context, tx sqlx.ExtContext, rendition Rendition) (Rendition, error) {
	rendition.Timestamps = db.Timestamps{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("renditions").
		Columns(
			"created_at",
			"updated_at",
			"photo_id",
			"original",
			"width",
			"height",
			"format",
			"rendition_configuration_id",
		).
		Values(
			rendition.CreatedAt,
			rendition.UpdatedAt,
			rendition.PhotoID,
			rendition.Original,
			rendition.Width,
			rendition.Height,
			rendition.Format,
			rendition.RenditionConfigurationID,
		).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return rendition, errors.Wrap(err, "could not create query")
	}

	err = tx.QueryRowxContext(ctx, sql, args...).Scan(&rendition.ID)
	if err != nil {
		return rendition, errors.Wrap(err, "could insert row")
	}

	return rendition, nil
}
