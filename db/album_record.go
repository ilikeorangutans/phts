package db

import (
	"log"
	"time"

	"github.com/ilikeorangutans/phts/pkg/database"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type AlbumRecord struct {
	Record
	Timestamps

	Name         string `db:"name" json:"name"`
	Slug         string `db:"slug" json:"slug"`
	CollectionID int64  `db:"collection_id" json:"collectionID"`
	PhotoCount   int    `db:"photo_count" json:"photoCount"`
	CoverPhotoID *int64 `db:"cover_photo_id" json:"coverPhotoID"`
}

type AlbumDB interface {
	FindByID(collectionID int64, id int64) (AlbumRecord, error)
	FindBySlug(collectionID int64, slug string) (AlbumRecord, error)
	List(collectionID int64, paginator database.Paginator) ([]AlbumRecord, error)
	Save(record AlbumRecord) (AlbumRecord, error)
	AddPhotos(collectionID int64, id int64, photoIDs []int64) error
	Delete(collectionID int64, id int64) error
}

func NewAlbumDB(db DB) AlbumDB {
	return &albumSQLDB{
		db:    db,
		clock: time.Now,
		sql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type albumSQLDB struct {
	db    DB
	clock Clock
	sql   sq.StatementBuilderType
}

func (a *albumSQLDB) albumsInCollection(collectionID int64) sq.SelectBuilder {
	return a.sql.
		Select("albums.*").
		From("albums").
		Where(
			sq.Eq{
				"collection_id": collectionID,
			},
		)
}

func (a *albumSQLDB) FindByID(collectionID int64, id int64) (AlbumRecord, error) {
	sql, args, _ := a.albumsInCollection(collectionID).Where(sq.Eq{"id": id}).ToSql()
	record := AlbumRecord{}
	err := a.db.QueryRowx(sql, args...).StructScan(&record)
	return record, err
}

func (a *albumSQLDB) FindBySlug(collectionID int64, slug string) (AlbumRecord, error) {
	sql, args, _ := a.albumsInCollection(collectionID).Where(sq.Eq{"slug": slug}).Limit(1).ToSql()
	record := AlbumRecord{}
	err := a.db.QueryRowx(sql, args...).StructScan(&record)
	return record, err
}

func (a *albumSQLDB) List(collectionID int64, paginator database.Paginator) ([]AlbumRecord, error) {
	sql, args, _ := paginator.Paginate(a.albumsInCollection(collectionID)).ToSql()
	result := []AlbumRecord{}
	err := a.db.Select(&result, sql, args...)
	return result, err
}

func (a *albumSQLDB) Save(record AlbumRecord) (AlbumRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(a.clock)
		sql, args, _ := a.sql.Update("albums").
			Set("name", record.Name).
			Set("cover_photo_id", record.CoverPhotoID).
			Set("updated_at", record.UpdatedAt.UTC()).
			Where(sq.Eq{
				"id": record.ID,
			}).
			ToSql()
		err = checkResult(a.db.Exec(sql, args...))
	} else {
		record.Timestamps = JustCreated(a.clock)
		sql, args, _ := a.sql.Insert("albums").
			Columns("name", "slug", "collection_id", "cover_photo_id", "created_at", "updated_at").
			Values(
				record.Name,
				record.Slug,
				record.CollectionID,
				record.CoverPhotoID,
				record.CreatedAt.UTC(),
				record.UpdatedAt.UTC(),
			).
			Suffix("RETURNING id").
			ToSql()
		err = a.db.QueryRow(sql, args...).Scan(&record.ID)
	}

	return record, err
}

func (a *albumSQLDB) AddPhotos(collectionID int64, id int64, photoIDs []int64) error {
	// TODO batching would probably be better here

	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	for _, photoID := range photoIDs {
		sql := "INSERT INTO album_photos (photo_id, album_id, created_at, updated_at) VALUES ($1, $2, $3, $4)"

		_, err = a.db.Exec(sql, photoID, id, a.clock().UTC(), a.clock().UTC())
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("error rolling back: %s", rollbackErr)
			}
			return err
		}
	}

	sql := "UPDATE albums SET photo_count = (SELECT COUNT(*) FROM album_photos WHERE album_id = $1) WHERE id = $1"
	_, err = a.db.Exec(sql, id)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("error rolling back: %s", rollbackErr)
		}
	}

	return tx.Commit()
}

func (a *albumSQLDB) Delete(collectionID int64, id int64) error {
	sql := "DELETE FROM albums WHERE collection_id = $1 AND id = $2"
	_, err := a.db.Exec(sql, collectionID, id)
	return err
}
