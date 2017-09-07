package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type PhotoRecord struct {
	Record
	Timestamps

	CollectionID   int64      `db:"collection_id"`
	RenditionCount int        `db:"rendition_count"`
	Description    *string    `db:"description"`
	Filename       string     `db:"filename"`
	TakenAt        *time.Time `db:"taken_at"`
}

type PhotoDB interface {
	FindByID(collectionID, id int64) (PhotoRecord, error)
	Save(record PhotoRecord) (PhotoRecord, error)
	List(collectionID int64, paginator Paginator) ([]PhotoRecord, error)
	Delete(collectionID, photoID int64) ([]int64, error)
}

func NewPhotoDB(db *sqlx.DB) PhotoDB {
	return &photoSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type photoSQLDB struct {
	db    *sqlx.DB
	clock Clock
}

type PhotoAndRendition struct {
	Photo     PhotoRecord
	Rendition RenditionRecord
}

func (c *photoSQLDB) Delete(collectionID, photoID int64) ([]int64, error) {
	result, err := c.db.Exec("DELETE FROM exif WHERE photo_id = $1", photoID)
	if err != nil {
		return nil, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	log.Printf("Removed %d exif records", count)

	rows, err := c.db.Queryx("DELETE FROM renditions WHERE photo_id = $1 RETURNING id", photoID)
	if err != nil {
		return nil, err
	}
	ids := []int64{}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	log.Printf("Removed %d rendition records", len(ids))

	result, err = c.db.Exec("DELETE FROM photos WHERE id = $1 and collection_id = $2", photoID, collectionID)
	if err != nil {
		return nil, err
	}
	count, err = result.RowsAffected()
	if err != nil {
		return nil, err
	}
	log.Printf("Removed %d photo records", count)

	return ids, nil
}

func (c *photoSQLDB) List(collection_id int64, paginator Paginator) ([]PhotoRecord, error) {
	sql, fields := paginator.Paginate("SELECT * FROM photos WHERE collection_id = $1", collection_id)
	log.Printf("Listing photos %s", sql)
	log.Printf("Fields %v", fields)
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
	if record.IsPersisted() {
		record.JustUpdated(c.clock)
		sql := "UPDATE photos SET collection_id = $1, filename = $2, rendition_count = $3, updated_at = $4 where id = $5"
		record.UpdatedAt = c.clock()
		err = checkResult(c.db.Exec(sql, record.CollectionID, record.Filename, record.RenditionCount, record.UpdatedAt.UTC(), record.ID))
	} else {
		record.Timestamps = JustCreated(c.clock)
		sql := "INSERT INTO photos (collection_id, filename, taken_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		log.Println(sql)
		err = c.db.QueryRow(sql, record.CollectionID, record.Filename, record.TakenAt, record.CreatedAt.UTC(), record.UpdatedAt.UTC()).Scan(&record.ID)
	}

	return record, err
}
