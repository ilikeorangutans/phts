package db

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type RenditionRecord struct {
	Record
	Timestamps

	PhotoID                  int64  `db:"photo_id"`
	Original                 bool   `db:"original"`
	Width                    uint   `db:"width"`
	Height                   uint   `db:"height"`
	Format                   string `db:"format"`
	RenditionConfigurationID int64  `db:"rendition_configuration_id"`
}

type RenditionDB interface {
	FindByID(id int64) (RenditionRecord, error)
	Save(RenditionRecord) (RenditionRecord, error)
	// TOOD FindBySize should rally be FindByRenditionConfiguration
	FindBySize(photoIDs []int64, width, height int) (map[int64]RenditionRecord, error)
	FindByRenditionConfiguration(photoIDs []int64, renditionConfigurationID int64) (map[int64]RenditionRecord, error)
	FindAllForPhoto(photoID int64) ([]RenditionRecord, error)
	// Create a new instace of RenditionRecord for th e given photo, filename, and binary data.
	Create(photo PhotoRecord, filename string, data []byte) (RenditionRecord, error)
}

func NewRenditionDB(db *sqlx.DB) RenditionDB {
	return &renditionSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type renditionSQLDB struct {
	db    *sqlx.DB
	clock func() time.Time
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
func (c *renditionSQLDB) FindByRenditionConfiguration(photoIDs []int64, renditionConfigurationID int64) (map[int64]RenditionRecord, error) {
	inQuery := []string{}
	for _, id := range photoIDs {
		inQuery = append(inQuery, fmt.Sprintf("%d", id))
	}

	sql := fmt.Sprintf("SELECT * FROM renditions WHERE rendition_configuration_id = $1 AND photo_id IN (%s) LIMIT $2", strings.Join(inQuery, ","))
	log.Printf("Executing  %s", sql)
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
	log.Printf("Executing  %s", sql)
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

func (c *renditionSQLDB) FindByID(id int64) (RenditionRecord, error) {
	var record RenditionRecord
	err := c.db.QueryRowx("SELECT * FROM renditions WHERE id = $1", id).StructScan(&record)
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
		sql := "INSERT INTO renditions (photo_id, original, width, height, format, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
		err = c.db.QueryRow(sql, record.PhotoID, record.Original, record.Width, record.Height, record.Format, record.CreatedAt.UTC(), record.UpdatedAt.UTC()).Scan(&record.ID)
	}

	return record, err
}
