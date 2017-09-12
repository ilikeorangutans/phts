package db

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/tiff"
)

const (
	exifTimeLayout = "2006:01:02 15:04:05"
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

// TODO this function belongs in a different package
func ExifRecordFromTiffTag(name string, tag *tiff.Tag) (ExifRecord, error) {
	record := ExifRecord{
		Type: int(tag.Type),
		Tag:  string(name),
	}

	if tag.Count > 1 {
		log.Printf("More than 1 value for %s of type %d: %d", name, tag.Type, tag.Count)
	}
	switch tag.Type {
	case tiff.DTByte, tiff.DTShort, tiff.DTLong, tiff.DTSShort, tiff.DTSLong:
		if num, err := tag.Int(0); err != nil {
			return record, nil
		} else {
			record.Num = int64(num)
		}
	case tiff.DTAscii:
		s, err := tag.StringVal()
		if err != nil {
			return record, err
		} else {
			record.StringValue = strings.TrimRight(s, "\x00")
			// TODO sanitize input values

			datetime, err := time.Parse(exifTimeLayout, record.StringValue)
			log.Printf("Parsed datetime: %v, %s", datetime, err)
			if err == nil {
				record.DateTime = &datetime
				return record, nil
			}

			if len(record.StringValue) == 0 {
				return record, fmt.Errorf("Skipping empty tag")
			}
		}
	case tiff.DTRational, tiff.DTSRational:
		if num, den, err := tag.Rat2(0); err != nil {
			return record, err
		} else {
			record.Num = num
			record.Denominator = den
		}
	case tiff.DTSByte:
	case tiff.DTUndefined:
	case tiff.DTFloat, tiff.DTDouble:
		f, err := tag.Float(0)
		if err != nil {
			return record, nil
		}

		record.Floating = f
	}
	return record, nil
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
