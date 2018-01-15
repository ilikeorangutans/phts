package db

import (
	"fmt"
	"time"
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
	List(collectionID int64, paginator Paginator) ([]AlbumRecord, error)
	Save(record AlbumRecord) (AlbumRecord, error)
}

func NewAlbumDB(db DB) AlbumDB {
	return &albumSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type albumSQLDB struct {
	db    DB
	clock Clock
}

func (a *albumSQLDB) FindByID(collectionID int64, id int64) (AlbumRecord, error) {
	return AlbumRecord{}, nil
}

func (a *albumSQLDB) FindBySlug(collectionID int64, slug string) (AlbumRecord, error) {
	return AlbumRecord{}, nil
}

func (a *albumSQLDB) List(collectionID int64, paginator Paginator) ([]AlbumRecord, error) {
	sql, fields := paginator.Paginate("SELECT * FROM albums WHERE collection_id = $1", collectionID)
	rows, err := a.db.Queryx(sql, fields...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []AlbumRecord{}
	for rows.Next() {
		record := AlbumRecord{}
		err = rows.StructScan(&record)
		if err != nil {
			return nil, err
		}

		result = append(result, record)
	}
	return result, nil
}

func (a *albumSQLDB) Save(record AlbumRecord) (AlbumRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(a.clock)
		err = fmt.Errorf("Implement me")
		//sql := "UPDATE photos SET filename = $1, updated_at = $2, rendition_count = (SELECT count(*) FROM renditions WHERE photo_id = $3) WHERE id = $3 AND collection_id = $4"
		//err = checkResult(c.db.Exec(sql, record.Filename, record.UpdatedAt.UTC(), record.ID, record.CollectionID))
	} else {
		record.Timestamps = JustCreated(a.clock)
		sql := "INSERT INTO albums (name, slug, collection_id, cover_photo_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
		err = a.db.QueryRow(
			sql,
			record.Name,
			record.Slug,
			record.CollectionID,
			record.CoverPhotoID,
			record.CreatedAt.UTC(),
			record.UpdatedAt.UTC(),
		).Scan(&record.ID)
	}

	return record, err
}
