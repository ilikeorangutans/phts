package db

import (
	"bytes"
	"fmt"
	"image"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type RenditionDB interface {
	FindByID(collectionID, id int64) (RenditionRecord, error)
	FindByPhotoAndConfigs(collectionID int64, photoID int64, renditionConfigurationIDs []int64) ([]RenditionRecord, error)
	Save(RenditionRecord) (RenditionRecord, error)
	// FindByRenditionConfiguration returns a map of photo id to rendition record
	FindByRenditionConfiguration(photoIDs []int64, renditionConfigurationID int64) (map[int64]RenditionRecord, error)
	// FindByRenditionConfigurations returns a map of photo ids to a set of rendition records matching the passed in ids
	FindByRenditionConfigurations(photoIDs []int64, renditionConfigurationIDs []int64) (map[int64][]RenditionRecord, error)
	FindAllForPhoto(photoID int64) ([]RenditionRecord, error)
	// Create a new instace of RenditionRecord for th e given photo, filename, and binary data.
	Create(photo PhotoRecord, filename string, data []byte) (RenditionRecord, error)
	DeleteForPhoto(photoID int64) ([]int64, error)
	FindByShareAndID(shareID, id int64) (RenditionRecord, error)
}

func NewRenditionDB(db DB) RenditionDB {
	return &renditionSQLDB{
		db:    db,
		clock: time.Now,
		sql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type renditionSQLDB struct {
	db    DB
	clock func() time.Time
	sql   sq.StatementBuilderType
}

func (c *renditionSQLDB) FindByShareAndID(shareID, id int64) (record RenditionRecord, err error) {
	sql, args, _ := c.sql.Select("r.*").
		From("renditions as r").
		Join("rendition_configurations rc on r.rendition_configuration_id = rc.id").
		Join("share_rendition_configurations src on rc.id = src.rendition_configuration_id").
		Where(sq.Eq{"src.share_id": shareID, "r.id": id}).
		ToSql()

	err = c.db.QueryRowx(sql, args...).StructScan(&record)

	return record, err
}

func (c *renditionSQLDB) DeleteForPhoto(photoID int64) ([]int64, error) {
	sql, args, _ := c.sql.Delete("renditions").Where(sq.Eq{"photo_id": photoID}).Suffix("RETURNING id").ToSql()

	var ids []int64
	rows, err := c.db.Queryx(sql, args)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}

func (c *renditionSQLDB) Create(photo PhotoRecord, filename string, data []byte) (RenditionRecord, error) {
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return RenditionRecord{}, err
	}

	record := RenditionRecord{
		PhotoID: photo.ID,
		Format:  format,
		Width:   uint(config.Width),
		Height:  uint(config.Height),
	}

	return record, nil
}

func (c *renditionSQLDB) FindAllForPhoto(photoID int64) ([]RenditionRecord, error) {
	sql, args, _ := c.sql.Select("renditions.*").
		From("renditions").
		Where(sq.Eq{"photo_id": photoID}).
		ToSql()

	result := []RenditionRecord{}
	err := c.db.Select(&result, sql, args...)

	return result, err
}

func (c *renditionSQLDB) FindByRenditionConfigurations(photoIDs []int64, renditionConfigurationIDs []int64) (map[int64][]RenditionRecord, error) {
	sql, args, err := c.sql.Select("renditions.*").
		From("renditions").
		Where(sq.Eq{
			"photo_id":                   photoIDs,
			"rendition_configuration_id": renditionConfigurationIDs,
		}).
		Limit(uint64(len(photoIDs) * len(renditionConfigurationIDs))).
		ToSql()

	if err != nil {
		return nil, err
	}

	var records []RenditionRecord
	c.db.Select(&records, sql, args...)

	result := make(map[int64][]RenditionRecord)
	for _, record := range records {
		result[record.PhotoID] = append(result[record.PhotoID], record)
	}

	return result, nil
}

func (c *renditionSQLDB) FindByRenditionConfiguration(photoIDs []int64, renditionConfigurationID int64) (map[int64]RenditionRecord, error) {
	sql, args, _ := c.sql.Select("renditions.*").
		From("renditions").
		Where(sq.Eq{
			"rendition_configuration_id": renditionConfigurationID,
			"photo_id":                   photoIDs,
		}).
		Limit(uint64(len(photoIDs))).
		ToSql()

	rows, err := c.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]RenditionRecord)
	for rows.Next() {
		rendition := RenditionRecord{}
		err = rows.StructScan(&rendition)
		if err != nil {
			return nil, err
		}

		result[rendition.PhotoID] = rendition
	}

	return result, nil
}

func (c *renditionSQLDB) FindByPhotoAndConfigs(collectionID int64, photoID int64, renditionConfigurationIDs []int64) ([]RenditionRecord, error) {
	var records []RenditionRecord

	idCount := len(renditionConfigurationIDs)
	// sqlx.IN doesn't take into account the number of existing params, so we'll explicitly setting the number of our param here:
	sql := fmt.Sprintf("SELECT r.* FROM renditions r, photos p WHERE p.id = r.photo_id AND r.rendition_configuration_id IN (?) AND p.id = $%d AND p.collection_id = $%d", idCount+1, idCount+2)
	query, args, err := sqlx.In(sql, renditionConfigurationIDs)
	if err != nil {
		return nil, err
	}
	query = c.db.Rebind(query)
	args = append(args, photoID)
	args = append(args, collectionID)
	records = []RenditionRecord{}
	err = c.db.Select(&records, query, args...)
	return records, err
}

func (c *renditionSQLDB) FindByID(collectionID, id int64) (RenditionRecord, error) {
	var record RenditionRecord
	sql := "SELECT r.* FROM renditions r, photos p WHERE r.id = $1 AND p.id = r.photo_id AND p.collection_id = $2"
	err := c.db.QueryRowx(sql, id, collectionID).StructScan(&record)
	return record, err
}

func (c *renditionSQLDB) Save(record RenditionRecord) (RenditionRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(c.clock)
		sql := "UPDATE renditions SET photo_id = $1, updated_at = $2 where id = $3"
		err = checkResult(c.db.Exec(sql, record.PhotoID, record.UpdatedAt.UTC(), record.ID))
	} else {
		record.Timestamps = JustCreated(c.clock)
		sql := "INSERT INTO renditions (photo_id, original, width, height, format, rendition_configuration_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
		err = c.db.QueryRow(
			sql,
			record.PhotoID,
			record.Original,
			record.Width,
			record.Height,
			record.Format,
			record.RenditionConfigurationID,
			record.CreatedAt.UTC(),
			record.UpdatedAt.UTC(),
		).Scan(&record.ID)
	}

	return record, err
}
