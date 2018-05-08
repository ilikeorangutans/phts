package db

import (
	"fmt"
	"time"

	sq "gopkg.in/Masterminds/squirrel.v1"
)

type PhotoDB interface {
	FindByID(collectionID, id int64) (PhotoRecord, error)
	Save(record PhotoRecord) (PhotoRecord, error)
	List(collectionID int64, paginator Paginator) ([]PhotoRecord, error)
	ListAlbum(collectionID int64, albumID int64, paginator Paginator) ([]PhotoRecord, error)
	Delete(collectionID, photoID int64) error
}

func NewPhotoDB(db DB) PhotoDB {
	return NewPhotoDBWithClock(db, time.Now)
}

func NewPhotoDBWithClock(db DB, clock Clock) PhotoDB {
	return &photoSQLDB{
		db:    db,
		clock: clock,
		sql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type photoSQLDB struct {
	db    DB
	clock Clock
	sql   sq.StatementBuilderType
}

func (c *photoSQLDB) Delete(collectionID, photoID int64) error {
	sql, args, _ := c.sql.Delete("photos").Where(sq.Eq{"collection_id": collectionID, "id": photoID}).ToSql()
	_, err := c.db.Exec(sql, args...)
	return err
}

func (c *photoSQLDB) photosInCollection(collectionID int64) sq.SelectBuilder {
	return c.sql.
		Select("photos.*").
		From("photos").
		Where(sq.Eq{
			"collection_id": collectionID,
		})
}

func (c *photoSQLDB) ListAlbum(collectionID int64, albumID int64, paginator Paginator) ([]PhotoRecord, error) {
	q := c.photosInCollection(collectionID).
		Join("album_photos on (photos.id = album_photos.photo_id)").
		Where(
			sq.Eq{
				"album_photos.album_id": albumID,
			},
		)

	paginator.ColumnPrefix = "photos"
	q = paginator.Paginate(q)
	sql, args, _ := q.ToSql()

	result := []PhotoRecord{}
	err := c.db.Select(&result, sql, args...)
	return result, err
}

func (c *photoSQLDB) List(collectionID int64, paginator Paginator) ([]PhotoRecord, error) {
	paginator.ColumnPrefix = "photos"
	q := c.photosInCollection(collectionID)
	q = paginator.Paginate(q)
	sql, args, _ := q.ToSql()

	result := []PhotoRecord{}
	err := c.db.Select(&result, sql, args...)
	return result, err
}

func (c *photoSQLDB) FindByID(collectionID, id int64) (PhotoRecord, error) {
	sql, args, _ := c.photosInCollection(collectionID).Where(sq.Eq{"id": id}).Limit(1).ToSql()
	var record PhotoRecord
	err := c.db.QueryRowx(sql, args...).StructScan(&record)
	return record, err
}

func (c *photoSQLDB) Save(record PhotoRecord) (PhotoRecord, error) {
	var err error
	if record.CollectionID < 1 {
		return record, fmt.Errorf("no collection id set")
	}

	if record.IsPersisted() {
		record.JustUpdated(c.clock)
		sql, args, _ := c.sql.
			Update("photos").
			Set("filename", record.Filename).
			Set("updated_at", record.UpdatedAt.UTC()).
			Where(sq.Eq{
				"id":            record.ID,
				"collection_id": record.CollectionID,
			}).
			ToSql()

		err = checkResult(c.db.Exec(sql, args...))
	} else {
		record.Timestamps = JustCreated(c.clock)
		sql, args, _ := c.sql.
			Insert("photos").
			Columns(
				"collection_id",
				"filename",
				"taken_at",
				"created_at",
				"updated_at",
			).
			Values(
				record.CollectionID,
				record.Filename,
				record.TakenAt,
				record.CreatedAt.UTC(),
				record.UpdatedAt.UTC(),
			).
			Suffix("RETURNING id").
			ToSql()

		err = c.db.QueryRow(
			sql,
			args...,
		).Scan(&record.ID)
	}

	return record, err
}
