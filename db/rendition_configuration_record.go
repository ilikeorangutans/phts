package db

import (
	"time"
)

type RenditionConfigurationRecord struct {
	Record
	Timestamps
	Width        int    `db:"width" json:"width"`
	Height       int    `db:"height" json:"height"`
	Name         string `db:"name" json:"name"`
	Quality      int    `db:"quality" json:"quality"`
	Private      bool   `db:"private" json:"private"`
	Resize       bool   `db:"resize" json:"resize"`
	CollectionID *int64 `db:"collection_id" json:"collection_id"`
}

func (r RenditionConfigurationRecord) Area() int64 {
	return int64(r.Width * r.Height)
}

type RenditionConfigurationDB interface {
	FindByID(int64, int64) (RenditionConfigurationRecord, error)
	FindByName(int64, string) (RenditionConfigurationRecord, error)
	Save(RenditionConfigurationRecord) (RenditionConfigurationRecord, error)
	FindForCollection(collectionID int64) ([]RenditionConfigurationRecord, error)
	Delete(int64) error
}

func NewRenditionConfigurationDB(db DB) RenditionConfigurationDB {
	return &renditionConfigurationSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type renditionConfigurationSQLDB struct {
	db    DB
	clock func() time.Time
}

func (c *renditionConfigurationSQLDB) FindByID(collectionID, id int64) (RenditionConfigurationRecord, error) {
	config := RenditionConfigurationRecord{}
	sql := "SELECT * FROM rendition_configurations WHERE (collection_id = $1 OR collection_id IS NULL) AND id = $2"
	err := c.db.QueryRowx(sql, collectionID, id).StructScan(&config)
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
			record.UpdatedAt.UTC(),
			record.CreatedAt.UTC(),
		).Scan(&record.ID)
	}

	return record, err
}

func (c *renditionConfigurationSQLDB) FindForCollection(collectionID int64) ([]RenditionConfigurationRecord, error) {
	sql := "SELECT * from rendition_configurations WHERE collection_id = $1 OR collection_id IS NULL ORDER BY width DESC, height DESC"
	rows, err := c.db.Queryx(sql, collectionID)
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
