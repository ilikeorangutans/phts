package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type CollectionRecord struct {
	Record
	Timestamps
	Sluggable
	Name string `db:"name"`
}

type CollectionDB interface {
	FindByID(id int64) (CollectionRecord, error)
	FindBySlug(slug string) (CollectionRecord, error)
	Save(collection CollectionRecord) (CollectionRecord, error)
}

type collectionSQLDB struct {
	db    *sqlx.DB
	clock func() time.Time
}

func (c *collectionSQLDB) FindByID(id int64) (CollectionRecord, error) {
	var record CollectionRecord
	err := c.db.QueryRow("SELECT * FROM collections WHERE id = $1", id).Scan(&record)
	return record, err
}

func (c *collectionSQLDB) FindBySlug(slug string) (CollectionRecord, error) {
	var record CollectionRecord
	err := c.db.QueryRow("SELECT * FROM collections WHERE slug = $1", slug).Scan(&record)
	return record, err
}

func checkResult(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	if count, err := result.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("No row updated")
	}
	return nil
}

func (c *collectionSQLDB) Save(record CollectionRecord) (CollectionRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated()
		sql := "UPDATE collections SET name = $1, slug = $2, updated_at = $3 where id = $4"
		record.UpdatedAt = c.clock()
		err = checkResult(c.db.Exec(sql, record.Name, record.Slug, record.UpdatedAt, record.ID))
	} else {
		record.Timestamps = JustCreated()
		sql := "INSERT INTO collections (name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
		err = c.db.QueryRow(sql, record.Name, record.Slug, record.CreatedAt, record.UpdatedAt).Scan(&record.ID)
	}

	return record, err
}
