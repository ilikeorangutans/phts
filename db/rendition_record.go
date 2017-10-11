package db

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type RenditionRecord struct {
	Record
	Timestamps

	PhotoID                  int64  `db:"photo_id" json:"photoID"`
	Original                 bool   `db:"original" json:"original"`
	Width                    uint   `db:"width" json:"width"`
	Height                   uint   `db:"height" json:"height"`
	Format                   string `db:"format" json:"format"`
	RenditionConfigurationID int64  `db:"rendition_configuration_id" json:"renditionConfigurationID"`
}

type RenditionDB interface {
	FindByID(collectionID, id int64) (RenditionRecord, error)
	FindByPhotoAndConfigs(collectionID int64, photoID int64, renditionConfigurationIDs []int64) ([]RenditionRecord, error)
	Save(RenditionRecord) (RenditionRecord, error)
	// TOOD FindBySize should rally be FindByRenditionConfiguration
	FindBySize(photoIDs []int64, width, height int) (map[int64]RenditionRecord, error)
	// FindByRenditionConfiguration returns a map of photo id to rendition record
	FindByRenditionConfiguration(photoIDs []int64, renditionConfigurationID int64) (map[int64]RenditionRecord, error)
	// FindByRenditionConfigurations returns a map of photo ids to a set of rendition records matching the passed in ids
	FindByRenditionConfigurations(photoIDs []int64, renditionConfigurationIDs []int64) (map[int64][]RenditionRecord, error)
	FindAllForPhoto(photoID int64) ([]RenditionRecord, error)
	// Create a new instace of RenditionRecord for th e given photo, filename, and binary data.
	Create(photo PhotoRecord, filename string, data []byte) (RenditionRecord, error)
	DeleteForPhoto(photoID int64) ([]int64, error)
}

func NewRenditionDB(db DB) RenditionDB {
	return &renditionSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type renditionSQLDB struct {
	db    DB
	clock func() time.Time
}

func (c *renditionSQLDB) DeleteForPhoto(photoID int64) ([]int64, error) {
	sql := "DELETE FROM RENDITIONS WHERE photo_id = $1 RETURNING id"

	var ids []int64
	rows, err := c.db.Queryx(sql, photoID)
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
	sql := "SELECT * FROM renditions WHERE photo_id = $1"
	rows, err := c.db.Queryx(sql, photoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []RenditionRecord{}
	for rows.Next() {
		record := RenditionRecord{}
		rows.StructScan(&record)

		result = append(result, record)
	}

	return result, nil
}

func (c *renditionSQLDB) FindByRenditionConfigurations(photoIDs []int64, renditionConfigurationIDs []int64) (map[int64][]RenditionRecord, error) {
	inPhotoIDs := []string{}
	for _, id := range photoIDs {
		inPhotoIDs = append(inPhotoIDs, fmt.Sprintf("%d", id))
	}

	inConfigIDs := []string{}
	for _, id := range renditionConfigurationIDs {
		inConfigIDs = append(inConfigIDs, fmt.Sprintf("%d", id))
	}
	sql := fmt.Sprintf("SELECT * FROM renditions WHERE rendition_configuration_id IN (%s) AND photo_id IN (%s)", strings.Join(inConfigIDs, ","), strings.Join(inPhotoIDs, ","))
	rows, err := c.db.Queryx(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]RenditionRecord)
	for rows.Next() {
		rendition := RenditionRecord{}
		err = rows.StructScan(&rendition)
		if err != nil {
			return nil, err
		}

		result[rendition.PhotoID] = append(result[rendition.PhotoID], rendition)
	}

	return result, nil
}

func (c *renditionSQLDB) FindByRenditionConfiguration(photoIDs []int64, renditionConfigurationID int64) (map[int64]RenditionRecord, error) {
	inQuery := []string{}
	for _, id := range photoIDs {
		inQuery = append(inQuery, fmt.Sprintf("%d", id))
	}

	sql := fmt.Sprintf("SELECT * FROM renditions WHERE rendition_configuration_id = $1 AND photo_id IN (%s) LIMIT $2", strings.Join(inQuery, ","))
	rows, err := c.db.Queryx(sql, renditionConfigurationID, len(photoIDs))
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

func (c *renditionSQLDB) FindBySize(photoIDs []int64, width, height int) (map[int64]RenditionRecord, error) {

	sizeConstraintField := "width"
	sizeConstraint := width
	if width == 0 {
		sizeConstraintField = "height"
		sizeConstraint = height
	}

	inQuery := []string{}
	for _, id := range photoIDs {
		inQuery = append(inQuery, fmt.Sprintf("%d", id))
	}

	sql := fmt.Sprintf("SELECT * FROM renditions WHERE %s = $1 and photo_id in (%s)", sizeConstraintField, strings.Join(inQuery, ","))
	rows, err := c.db.Queryx(sql, sizeConstraint)
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
