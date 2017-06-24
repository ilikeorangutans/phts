package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type PhotoRecord struct {
	Record
	Timestamps

	CollectionID   int64   `db:"collection_id"`
	RenditionCount int     `db:"rendition_count"`
	Description    *string `db:"description"`
}

type PhotoDB interface {
	FindByID(id int64) (PhotoRecord, error)
	Save(record PhotoRecord) (PhotoRecord, error)
	List(collectionID int64, afterID int64, orderBy string, limit int) ([]PhotoRecord, error)
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

type PhotoAndRendition struct {
	Photo     PhotoRecord
	Rendition RenditionRecord
}

func (c *photoSQLDB) List(collection_id int64, afterID int64, orderBy string, limit int) ([]PhotoRecord, error) {
	sql := "SELECT * FROM photos WHERE collection_id = $1 AND id > $2 ORDER BY $3 LIMIT $4"

	rows, err := c.db.Queryx(sql, collection_id, afterID, orderBy, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []PhotoRecord{}
	for rows.Next() {
		record := PhotoRecord{}
		err = rows.StructScan(&record)
		if err != nil {
			return nil, err
		}

		result = append(result, record)
	}
	return result, nil
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
