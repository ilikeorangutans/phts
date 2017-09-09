package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type RenditionConfigurationRecord struct {
	Record
	Timestamps
	Width        int    `db:"width"`
	Height       int    `db:"height"`
	Name         string `db:"name"`
	Quality      int    `db:"quality"`
	CollectionID *int64 `db:"collection_id"`
}

type RenditionConfigurationDB interface {
	FindByID(int64, int64) (RenditionConfigurationRecord, error)
	FindByName(int64, string) (RenditionConfigurationRecord, error)
	Save(RenditionConfigurationRecord) (RenditionConfigurationRecord, error)
	FindForCollection(collectionID int64) ([]RenditionConfigurationRecord, error)
	Delete(int64) error
}

func NewRenditionConfigurationDB(db *sqlx.DB) RenditionConfigurationDB {
	return &renditionConfigurationSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type renditionConfigurationSQLDB struct {
	db    *sqlx.DB
	clock func() time.Time
}

func (c *renditionConfigurationSQLDB) FindByID(collectionID, id int64) (RenditionConfigurationRecord, error) {
	config := RenditionConfigurationRecord{}
	err := c.db.QueryRowx("SELECT * FROM rendition_configurations WHERE collection_id = $1 OR collection_id IS NULL AND id = $2 LIMIT 1", collectionID, id).StructScan(&config)
	return config, err
}

func (c *renditionConfigurationSQLDB) Delete(id int64) error {
	_, err := c.db.Exec(
		"DELETE FROM rendition_configurations WHERE id=$1",
		id,
	)

	return err
}

func (c *renditionConfigurationSQLDB) FindByName(collectionID int64, name string) (RenditionConfigurationRecord, error) {
	config := RenditionConfigurationRecord{}
	err := c.db.QueryRowx("SELECT * FROM rendition_configurations WHERE collection_id = $1 OR collection_id IS NULL AND name = $2 LIMIT 1", collectionID, name).StructScan(&config)
	return config, err
}

func (c *renditionConfigurationSQLDB) Save(record RenditionConfigurationRecord) (RenditionConfigurationRecord, error) {
	var err error

	if record.IsPersisted() {
		record.JustUpdated(c.clock)
		sql := "UPDATE rendition_configurations SET width=$1, height=$2, name=$3, quality=$4, updated_at=$5 WHERE collection_id=$6 AND id=$7"
		err = checkResult(c.db.Exec(
			sql,
			record.Width,
			record.Height,
			record.Name,
			record.Quality,
			record.UpdatedAt.UTC(),
			record.CollectionID,
			record.ID,
		))
	} else {
		record.Timestamps = JustCreated(c.clock)
		sql := "INSERT INTO rendition_configurations (width, height, name, quality, collection_id, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
		err = c.db.QueryRow(
			sql,
			record.Width,
			record.Height,
			record.Name,
			record.Quality,
			record.CollectionID,
			record.UpdatedAt,
			record.CreatedAt,
		).Scan(&record.ID)
	}

	return record, err
}

func (c *renditionConfigurationSQLDB) FindForCollection(collectionID int64) ([]RenditionConfigurationRecord, error) {
	rows, err := c.db.Queryx("SELECT * from rendition_configurations WHERE collection_id = $1 OR collection_id IS NULL", collectionID)
	if err != nil {
		return nil, err
	}

	result := []RenditionConfigurationRecord{}
	for rows.Next() {
		record := RenditionConfigurationRecord{}
		err := rows.StructScan(&record)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}

	return result, nil
}