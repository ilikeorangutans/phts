package model

import (
	"context"
	"log"

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

// FindOriginalRenditionByPhoto finds the original rendition for a photo
func FindOriginalRenditionByPhoto(ctx context.Context, tx sqlx.QueryerContext, photo Photo) (Rendition, error) {
	sql, args := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("renditions").
		Where(sq.Eq{
			"photo_id": photo.ID,
			"original": true,
		}).
		Limit(1).
		MustSql()

	var rendition Rendition
	err := sqlx.GetContext(ctx, tx, &rendition, sql, args...)
	if err != nil {
		return Rendition{}, errors.Wrap(err, "could get rendition")
	}

	return rendition, nil
}

// CreateRendition creates a new rendition in the database.
func InsertRendition(ctx context.Context, tx sqlx.ExtContext, rendition Rendition) (Rendition, error) {
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
		log.Printf("%+v", err)
		// TODO this query sometimes fails:
		// error adding rendition to photo: pq: duplicate key value violates unique constraint "renditions_photo_id_rendition_configuration_id_idx"
		// could insert row
		return rendition, errors.Wrap(err, "could not insert row")
	}

	return rendition, nil
}

func FindRenditionsForPhoto(ctx context.Context, tx sqlx.QueryerContext, photo Photo, renditionConfigurations ...RenditionConfiguration) ([]Rendition, error) {
	var configIDs []int64
	for _, config := range renditionConfigurations {
		configIDs = append(configIDs, config.ID)
	}
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("renditions").
		Where(sq.Eq{
			"photo_id":                   photo.ID,
			"rendition_configuration_id": configIDs,
		}).
		Limit(uint64(len(renditionConfigurations))).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not create query")
	}

	var renditions []Rendition
	err = sqlx.SelectContext(ctx, tx, &renditions, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch rows")
	}
	return renditions, nil
}
