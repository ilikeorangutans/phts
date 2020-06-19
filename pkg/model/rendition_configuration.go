package model

import (
	"context"
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

type RenditionConfiguration struct {
	db.Record
	db.Timestamps

	CollectionID *int64 `db:"collection_id" json:"collectionID"`
	Height       int    `db:"height" json:"height"`
	Name         string `db:"name" json:"name"`
	Original     bool   `db:"original" json:"original"`
	Private      bool   `db:"private" json:"private"` // TODO rename this to "system" or "reserved" or "locked"
	Quality      int    `db:"quality" json:"quality"`
	Resize       bool   `db:"resize" json:"resize"`
	Width        int    `db:"width" json:"width"`
}

func FindOriginalRenditionConfiguration(ctx context.Context, dbx sqlx.ExtContext) (RenditionConfiguration, error) {
	var config RenditionConfiguration
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("rendition_configurations").
		Where(sq.Eq{"original": true, "collection_id": nil}).
		Limit(1).
		ToSql()
	if err != nil {
		return config, errors.Wrap(err, "could not create query")
	}

	err = dbx.QueryRowxContext(ctx, sql, args...).StructScan(&config)
	if err != nil {
		return config, errors.Wrap(err, "could not query row")
	}

	return config, nil
}

// FindMissingRenditionConfigurations finds all rendition configurations for which the given photo doesn't have
// a rendition.
func FindMissingRenditionConfigurations(ctx context.Context, dbx sqlx.QueryerContext, photo Photo) ([]RenditionConfiguration, error) {
	sql := `
	  select
		*
	  from
		rendition_configurations
	  where
		(
		  collection_id is null
		  or
		  collection_id = $1
		)
		and
		id not in (
		  select
			rendition_configuration_id
		  from
			renditions
		  where
			photo_id = $2
		)
	`

	log.Printf("%s", sql)
	rows, err := dbx.QueryxContext(ctx, sql, photo.CollectionID, photo.ID)
	if err != nil {
		return nil, errors.Wrap(err, "could not query")
	}

	var configs []RenditionConfiguration
	for rows.Next() {
		var config RenditionConfiguration
		rows.StructScan(&config)
		configs = append(configs, config)
	}

	return configs, nil
}

// FindApplicableRenditionConfigurations finds all RenditionConfigurations that are applicable to the given
// collection, ignoring the original rendition.
func FindApplicableRenditionConfigurations(ctx context.Context, dbx sqlx.QueryerContext, collection Collection) ([]RenditionConfiguration, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("rendition_configurations").
		Where(
			sq.And{
				sq.Eq{"original": false},
				sq.Or{
					sq.Eq{
						"collection_id": collection.ID,
					},
					sq.Eq{
						"collection_id": nil,
					},
				},
			},
		).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not create query")
	}

	rows, err := dbx.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "could not query")
	}

	var configs []RenditionConfiguration
	for rows.Next() {
		var config RenditionConfiguration
		rows.StructScan(&config)
		configs = append(configs, config)
	}

	return configs, nil
}
