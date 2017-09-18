package db

import (
	"fmt"
	"log"
	"time"
)

type ExifRecord struct {
	Record
	Timestamps

	PhotoID     int64      `db:"photo_id"`
	Type        int        `db:"value_type"`
	Tag         string     `db:"tag"`
	StringValue string     `db:"string"`
	Num         int64      `db:"num"`
	Denominator int64      `db:"denom"`
	DateTime    *time.Time `db:"datetime"`
	Floating    float64    `db:"floating"`
}

type ExifDB interface {
	ByTag(photoID int64, tag string) (ExifRecord, error)
	AllForPhoto(photoID int64) ([]ExifRecord, error)
	Save(photoID int64, record ExifRecord) (ExifRecord, error)
}

func NewExifDB(db DB) ExifDB {
	return &exifSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type exifSQLDB struct {
	db    DB
	clock Clock
}

func (e *exifSQLDB) Save(photoID int64, record ExifRecord) (ExifRecord, error) {
	sql := "INSERT INTO exif (photo_id, value_type, tag, string, num, denom, datetime, floating) values ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	log.Printf("record: %v", record)
	log.Printf("Saving %s", sql)
	if record.IsPersisted() {
		// TODO do we ever update exif tags?
		return record, fmt.Errorf("exif tags cannot be updated")
	}
	record.Timestamps = JustCreated(e.clock)

	err := e.db.QueryRowx(sql, photoID, record.Type, record.Tag, record.StringValue, record.Num, record.Denominator, record.DateTime, record.Floating).Scan(&record.ID)

	return record, err
}

func (e *exifSQLDB) AllForPhoto(photoID int64) ([]ExifRecord, error) {
	sql := "SELECT * FROM exif WHERE photo_id = $1"

	rows, err := e.db.Queryx(sql, photoID)
	if err != nil {
		return nil, err
	}

	records := []ExifRecord{}
	for rows.Next() {
		record := ExifRecord{}
		err := rows.StructScan(&record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (e *exifSQLDB) ByTag(photoID int64, tag string) (ExifRecord, error) {
	record := ExifRecord{}
	sql := "SELECT * FROM exif WHERE photo_id = $1 and tag = $2"
	e.db.QueryRowx(sql, photoID, tag).StructScan(&record)

	return record, nil
}
