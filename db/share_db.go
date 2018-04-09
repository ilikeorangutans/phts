package db

import (
	"time"

	sq "gopkg.in/Masterminds/squirrel.v1"
)

type ShareDB interface {
	FindByShareSiteAndSlug(shareSiteID int64, slug string) (ShareRecord, error)
	FindByPhoto(photoID int64) ([]ShareRecord, error)
	Save(ShareRecord) (ShareRecord, error)
}

func NewShareDB(dbx DB) ShareDB {
	return &shareSQLDB{
		db:    dbx,
		clock: time.Now,
		sql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type shareSQLDB struct {
	db    DB
	clock Clock
	sql   sq.StatementBuilderType
}

func (c *shareSQLDB) Save(record ShareRecord) (ShareRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(c.clock)

		sql, args, _ := c.sql.Update("shares").
			Where(sq.Eq{"id": record.ID}).
			Set("slug", record.Slug).
			Set("updated_at", record.UpdatedAt.UTC()).
			ToSql()

		err = checkResult(c.db.Exec(sql, args...))
	} else {
		record.Timestamps = JustCreated(c.clock)

		sql, args, _ := c.sql.Insert("shares").
			Columns("photo_id", "collection_id", "share_site_id", "slug", "created_at", "updated_at").
			Values(record.PhotoID, record.CollectionID, record.ShareSiteID, record.Slug, record.CreatedAt.UTC(), record.UpdatedAt.UTC()).
			Suffix("RETURNING id").
			ToSql()

		err = c.db.QueryRow(sql, args...).Scan(&record.ID)
	}

	return record, err
}

func (c *shareSQLDB) FindByShareSiteAndSlug(shareSiteID int64, slug string) (ShareRecord, error) {
	sql, args, _ := c.sql.Select("shares.*").
		From("shares").
		Where(sq.Eq{
			"share_site_id": shareSiteID,
			"slug":          slug,
		}).
		Limit(1).
		ToSql()

	var record ShareRecord
	err := c.db.QueryRowx(sql, args...).StructScan(&record)
	return record, err
}

func (c *shareSQLDB) FindByPhoto(photoID int64) ([]ShareRecord, error) {
	sql, args, _ := c.sql.Select("shares.*").
		From("shares").
		Where(sq.Eq{"photo_id": photoID}).
		ToSql()

	var result []ShareRecord
	err := c.db.Select(&result, sql, args...)
	return result, err
}
