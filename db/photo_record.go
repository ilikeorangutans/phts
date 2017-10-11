package db

import (
	"fmt"
	"log"
	"time"
)

type PhotoRecord struct {
	Record
	Timestamps

	CollectionID   int64      `db:"collection_id" json:"collectionID"`
	RenditionCount int        `db:"rendition_count" json:"renditionCount"`
	Description    string     `db:"description" json:"description"`
	Filename       string     `db:"filename" json:"filename"`
	TakenAt        *time.Time `db:"taken_at" json:"takenAt"`
}

type PhotoDB interface {
	FindByID(collectionID, id int64) (PhotoRecord, error)
	Save(record PhotoRecord) (PhotoRecord, error)
	List(collectionID int64, paginator Paginator) ([]PhotoRecord, error)
	Delete(collectionID, photoID int64) error
}

func NewPhotoDB(db DB) PhotoDB {
	return &photoSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type photoSQLDB struct {
	db    DB
	clock Clock
}

type PhotoAndRendition struct {
	Photo     PhotoRecord
	Rendition RenditionRecord
}

func (c *photoSQLDB) Delete(collectionID, photoID int64) error {
	sql := "DELETE FROM photos WHERE id = $1 and collection_id = $2"
	result, err := c.db.Exec(sql, photoID, collectionID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("expected to delete one row, but none was deleted")
	}

	return nil
}

func (c *photoSQLDB) List(collectionID int64, paginator Paginator) ([]PhotoRecord, error) {
	sql, fields := paginator.Paginate("SELECT * FROM photos WHERE collection_id = $1", collectionID)
	rows, err := c.db.Queryx(sql, fields...)
	if err != nil {
		log.Panic(err)
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

func (c *photoSQLDB) FindByID(collectionID, id int64) (PhotoRecord, error) {
	var record PhotoRecord
	err := c.db.QueryRowx("SELECT * FROM photos WHERE collection_id = $1 AND id = $2", collectionID, id).StructScan(&record)
	return record, err
}

func (c *photoSQLDB) Save(record PhotoRecord) (PhotoRecord, error) {
	var err error
	if record.CollectionID < 1 {
		return record, fmt.Errorf("no collection id set")
	}

	if record.IsPersisted() {
		record.JustUpdated(c.clock)
		sql := "UPDATE photos SET filename = $1, updated_at = $2, rendition_count = (SELECT count(*) FROM renditions WHERE photo_id = $3) where id = $3 AND collection_id = $4"
		err = checkResult(c.db.Exec(sql, record.Filename, record.UpdatedAt.UTC(), record.ID, record.CollectionID))
	} else {
		record.Timestamps = JustCreated(c.clock)
		sql := "INSERT INTO photos (collection_id, filename, taken_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		err = c.db.QueryRow(sql, record.CollectionID, record.Filename, record.TakenAt, record.CreatedAt.UTC(), record.UpdatedAt.UTC()).Scan(&record.ID)
	}

	return record, err
}
