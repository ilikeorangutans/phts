package db

import (
	"fmt"
	"strings"
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
	fieldNames := []string{"rc.id", "rc.created_at", "rc.updated_at", "rc.width", "rc.height", "rc.name", "rc.quality", "rc.private", "rc.resize", "rc.original", "rc.collection_id"}
	var rcFields []string
	for _, name := range fieldNames {
		rcFields = append(rcFields, fmt.Sprintf("%s as \"%s\"", name, name))
	}
	sql, args, err := s.sql.Select(fmt.Sprintf("share_rendition_configurations.*, %s", strings.Join(rcFields, ","))).
		From("share_rendition_configurations").
		Join("rendition_configurations AS rc ON share_rendition_configurations.rendition_configuration_id = rc.id ").
		Where(sq.Eq{"share_id": shareID}).
		ToSql()

	var result []ShareRenditionConfigurationRecord
	err = s.db.Select(&result, sql, args...)
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
		query := s.sql.Insert("share_rendition_configurations").
			Columns("share_id", "rendition_configuration_id", "created_at", "updated_at")

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
