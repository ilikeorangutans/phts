package model

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ShareSite struct {
	db.Record
	db.Timestamps
	Domain string `db:"domain" json:"domain"`
}

func FindShareSiteByDomain(ctx context.Context, tx sqlx.QueryerContext, domain string) (ShareSite, error) {
	var shareSite ShareSite
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("share_sites").
		Where(sq.Eq{"domain": domain}).
		Limit(1).
		ToSql()
	if err != nil {
		return shareSite, errors.Wrap(err, "could not build query")
	}

	err = tx.QueryRowxContext(ctx, sql, args...).StructScan(&shareSite)
	if err != nil {
		return shareSite, errors.Wrap(err, "could not query")
	}

	return shareSite, nil
}
