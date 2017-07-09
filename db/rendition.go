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

	PhotoID  int64  `db:"photo_id"`
	Original bool   `db:"original"`
	Width    uint   `db:"width"`
	Height   uint   `db:"height"`
	Format   string // TODO: add to database
}

type RenditionConfiguration struct {
	Record
	Width        int       `db:"width"`
	Height       int       `db:"height"`
	Name         string    `db:"name"`
	Quality      int       `db:"quality"`
	CollectionID *int64    `db:"collection_id"`
	CreatedAt    time.Time `db:"created_at"`
}

func NewRenditionRecord(photo PhotoRecord, filename string, data []byte) (RenditionRecord, error) {
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return RenditionRecord{}, err
	}

	record := RenditionRecord{
		Timestamps: JustCreated(),
		PhotoID:    photo.ID,
		Format:     format,
		Width:      uint(config.Width),
		Height:     uint(config.Height),
	}

	return record, nil
}

type RenditionDB interface {
	FindByID(id int64) (RenditionRecord, error)
	Save(RenditionRecord) (RenditionRecord, error)
	FindBySize(photoIDs []int64, width, height int) (map[int64]RenditionRecord, error)
	FindAllForPhoto(photoID int64) ([]RenditionRecord, error)
	ApplicableConfigs(collectionID int64) ([]RenditionConfiguration, error)
	FindConfig(collectionID int64, name string) (RenditionConfiguration, error)
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

func (c *renditionSQLDB) FindConfig(collectionID int64, name string) (RenditionConfiguration, error) {
	config := RenditionConfiguration{}
	err := c.db.QueryRowx("SELECT * FROM rendition_configurations WHERE collection_id = $1 OR collection_id IS NULL AND name = $2", collectionID, name).StructScan(&config)
	return config, err
}

func (c *renditionSQLDB) ApplicableConfigs(collectionID int64) ([]RenditionConfiguration, error) {
	rows, err := c.db.Queryx("SELECT * from rendition_configurations WHERE collection_id = $1 OR collection_id IS NULL", collectionID)
	if err != nil {
		return nil, err
	}

	result := []RenditionConfiguration{}
	for rows.Next() {
		record := RenditionConfiguration{}
		err := rows.StructScan(&record)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}

	return result, nil
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

func (c *renditionSQLDB) FindBySize(photo_ids []int64, width, height int) (map[int64]RenditionRecord, error) {

	sizeConstraintField := "width"
	sizeConstraint := width
	if width == 0 {
		sizeConstraintField = "height"
		sizeConstraint = height
	}

	inQuery := []string{}
	for _, id := range photo_ids {
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
		record.JustUpdated()
		sql := "UPDATE renditions SET photo_id = $1, updated_at = $2 where id = $3"
		record.UpdatedAt = c.clock()
		err = checkResult(c.db.Exec(sql, record.PhotoID, record.UpdatedAt.UTC(), record.ID))
	} else {
		record.Timestamps = JustCreated()
		sql := "INSERT INTO renditions (photo_id, original, width, height, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
		err = c.db.QueryRow(sql, record.PhotoID, record.Original, record.Width, record.Height, record.CreatedAt.UTC(), record.UpdatedAt.UTC()).Scan(&record.ID)
	}

	return record, err
}
