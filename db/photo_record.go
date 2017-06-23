package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type PhotoRecord struct {
	Record
	Timestamps

	CollectionID int64 `db:"collection_id"`
}

type PhotoDB interface {
	FindByID(id int64) (PhotoRecord, error)
	Save(record PhotoRecord) (PhotoRecord, error)
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

func (c *photoSQLDB) FindByID(id int64) (PhotoRecord, error) {
	var record PhotoRecord
	err := c.db.QueryRow("SELECT * FROM photos WHERE id = $1", id).Scan(&record)
	return record, err
}

func (c *photoSQLDB) Save(record PhotoRecord) (PhotoRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated()
		sql := "UPDATE photos SET collection_id = $1, updated_at = $2 where id = $3"
		record.UpdatedAt = c.clock()
		err = checkResult(c.db.Exec(sql, record.CollectionID, record.UpdatedAt, record.ID))
	} else {
		record.Timestamps = JustCreated()
		sql := "INSERT INTO collections (collection_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id"
		err = c.db.QueryRow(sql, record.CollectionID, record.CreatedAt, record.UpdatedAt).Scan(&record.ID)
	}

	return record, err
}
