package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type PhotoRecord struct {
	Record
	Timestamps

	CollectionID   int64 `db:"collection_id"`
	RenditionCount int   `db:"rendition_count"`
}

type PhotoDB interface {
	FindByID(id int64) (PhotoRecord, error)
	Save(record PhotoRecord) (PhotoRecord, error)
	ListWithRenditions(int) error
}

func NewPhotoDB(db *sqlx.DB) PhotoDB {
	return &photoSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type photoSQLDB struct {
	db    *sqlx.DB
	clock func() time.Time
}

func (c *photoSQLDB) ListWithRenditions(count int) (PhotoRecord, error) {
	log.Printf("ListWithRenditions()")

	c.db.Queryx("SELECT ")

}

func (c *photoSQLDB) FindByID(id int64) (PhotoRecord, error) {
	var record PhotoRecord
	err := c.db.QueryRow("SELECT * FROM photos WHERE id = $1", id).Scan(&record)
	return record, err
}

func (c *photoSQLDB) Save(record PhotoRecord) (PhotoRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated()
		sql := "UPDATE photos SET collection_id = $1, rendition_count = $2, updated_at = $3 where id = $4"
		record.UpdatedAt = c.clock()
		err = checkResult(c.db.Exec(sql, record.CollectionID, record.RenditionCount, record.UpdatedAt, record.ID))
	} else {
		record.Timestamps = JustCreated()
		sql := "INSERT INTO photos (collection_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id"
		err = c.db.QueryRow(sql, record.CollectionID, record.CreatedAt, record.UpdatedAt).Scan(&record.ID)
	}

	return record, err
}
