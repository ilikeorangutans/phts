package db

import (
	"log"
	"time"

	sq "gopkg.in/Masterminds/squirrel.v1"
)

type ShareRenditionConfigurationDB interface {
	FindByShare(shareID int64) ([]ShareRenditionConfigurationRecord, error)
	SetForShare(shareID int64, configs []ShareRenditionConfigurationRecord) ([]ShareRenditionConfigurationRecord, error)
}

func NewShareRenditionConfigurationDB(dbx DB) ShareRenditionConfigurationDB {
	return &shareRenditionConfigurationSQLDB{
		db:    dbx,
		clock: time.Now,
		sql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type shareRenditionConfigurationSQLDB struct {
	db    DB
	clock Clock
	sql   sq.StatementBuilderType
}

func (s *shareRenditionConfigurationSQLDB) FindByShare(shareID int64) ([]ShareRenditionConfigurationRecord, error) {
	sql, args, _ := s.sql.Select("share_rendition_configurations.*").
		From("share_rendition_configurations").
		Where(sq.Eq{"share_id": shareID}).
		ToSql()

	var result []ShareRenditionConfigurationRecord
	err := s.db.Select(&result, sql, args...)
	return result, err
}

func (s *shareRenditionConfigurationSQLDB) SetForShare(shareID int64, configs []ShareRenditionConfigurationRecord) ([]ShareRenditionConfigurationRecord, error) {
	existingConfigs, err := s.FindByShare(shareID)
	if err != nil {
		return nil, err
	}

	have := make(map[int64]struct{})
	for _, config := range existingConfigs {
		have[config.RenditionConfigurationID] = struct{}{}
	}
	want := make(map[int64]struct{})
	for _, config := range configs {
		want[config.RenditionConfigurationID] = struct{}{}
	}
	_, add, remove := partitionIDs(want, have)

	if len(add) > 0 {
		query := s.sql.Insert("share_rendition_configurations").Columns("share_id", "rendition_configuration_id", "created_at", "updated_at")
		for _, id := range add {
			query = query.Values(shareID, id, s.clock().UTC(), s.clock().UTC())
		}

		sql, args, _ := query.ToSql()
		_, err := s.db.Exec(sql, args...)
		if err != nil {
			return nil, err
		}
	}
	if len(remove) > 0 {
		sql, args, _ := s.sql.
			Delete("share_rendition_configurations").
			Where(sq.Eq{"share_id": shareID, "rendition_configuration_id": remove}).
			ToSql()

		_, err := s.db.Exec(sql, args...)
		if err != nil {
			return nil, err
		}
	}
	return configs, nil
}

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

func partitionIDs(want, have map[int64]struct{}) (keep, add, remove []int64) {
	for id := range have {
		if _, ok := want[id]; ok {
		} else {
			remove = append(remove, id)
		}
	}

	for id := range want {
		add = append(add, id)
	}

	return keep, add, remove
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

	// TODO might be nicer to use .Select()
	rows, err := c.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	var result []ShareRecord
	for rows.Next() {
		record := ShareRecord{}
		err = rows.StructScan(&record)
		if err != nil {
			log.Printf("Error scanning: %s", err)
			return nil, err
		}
		result = append(result, record)
	}
	return result, err
}
