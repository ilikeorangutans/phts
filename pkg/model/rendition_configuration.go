package model

import (
	"context"

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

func FindNonOriginalRenditionConfigurations(ctx context.Context, dbx sqlx.QueryerContext, collection Collection) ([]RenditionConfiguration, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("rendition_configurations").
		Where(sq.Eq{"original": false, "collection_id": collection.ID}).
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
