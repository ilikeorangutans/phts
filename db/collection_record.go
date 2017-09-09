package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// CollectionRecord is a single database level record of a collection.
type CollectionRecord struct {
	Record
	Timestamps
	Sluggable
	Name       string `db:"name"`
	PhotoCount int    `db:"photo_count"`
}

type CollectionDB interface {
	FindByID(id int64) (CollectionRecord, error)
	FindBySlug(slug string) (CollectionRecord, error)
	Save(collection CollectionRecord) (CollectionRecord, error)
	List(count int, afterID int64, orderBy string) ([]CollectionRecord, error)
	Delete(int64) error
}

func NewCollectionDB(db *sqlx.DB) CollectionDB {
	return &collectionSQLDB{
		clock: time.Now,
		db:    db,
	}
}

type collectionSQLDB struct {
	db    *sqlx.DB
	clock Clock
}

func (c *collectionSQLDB) List(count int, afterID int64, orderBy string) ([]CollectionRecord, error) {
	result := []CollectionRecord{}
	rows, err := c.db.Queryx("SELECT * FROM collections WHERE id > $1 order by $2 limit $3", afterID, orderBy, count)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		record := CollectionRecord{}
		err := rows.StructScan(&record)
		if err != nil {
			return result, err
		}
		result = append(result, record)
	}

	return result, nil
}

func (c *collectionSQLDB) FindByID(id int64) (CollectionRecord, error) {
	var record CollectionRecord
	err := c.db.QueryRow("SELECT * FROM collections WHERE id = $1 LIMIT 1", id).Scan(&record)
	return record, err
}

func (c *collectionSQLDB) Delete(id int64) error {
	sql := "DELETE FROM collections WHERE id=$1"
	return checkResult(c.db.Exec(sql, id))
}

func (c *collectionSQLDB) FindBySlug(slug string) (CollectionRecord, error) {
	var record CollectionRecord
	err := c.db.QueryRowx("SELECT * FROM collections WHERE slug = $1 LIMIT 1", slug).StructScan(&record)
	return record, err
}

func (c *collectionSQLDB) Save(record CollectionRecord) (CollectionRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(c.clock)
		sql := "UPDATE collections SET name = $1, slug = $2, updated_at = $3, photo_count = (SELECT count(*) FROM photos WHERE collection_id = $4) WHERE id = $4"
		record.UpdatedAt = c.clock()
		err = checkResult(c.db.Exec(
			sql,
			record.Name,
			record.Slug,
			record.UpdatedAt,
			record.ID,
		))
	} else {
		record.Timestamps = JustCreated(c.clock)
		sql := "INSERT INTO collections (name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
		err = c.db.QueryRow(sql, record.Name, record.Slug, record.CreatedAt, record.UpdatedAt).Scan(&record.ID)
	}

	return record, err
}
