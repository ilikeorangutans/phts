package db

import (
	"fmt"
	"time"

	"github.com/rwcarlsen/goexif/tiff"
)

type ExifRecord struct {
	Record
	Timestamps `json:"-"`

	PhotoID     int64      `db:"photo_id" json:"-"`
	Type        uint16     `db:"value_type" json:"-"`
	TypeName    string     `db:"-" json:"type"`
	Tag         string     `db:"tag" json:"tag"`
	StringValue string     `db:"string" json:"string_value"`
	Num         int64      `db:"num" json:"number"`
	Denominator int64      `db:"denom" json:"denominator"`
	DateTime    *time.Time `db:"datetime" json:"datetime"`
	Floating    float64    `db:"floating" json:"float"`
}

func (e ExifRecord) String() string {
	switch tiff.DataType(e.Type) {
	case tiff.DTAscii:
		return fmt.Sprintf("%s", e.StringValue)
	case tiff.DTLong, tiff.DTShort:
		return fmt.Sprintf("%d", e.Num)
	case tiff.DTRational, tiff.DTSRational:
		return fmt.Sprintf("%d/%d", e.Num, e.Denominator)
	}
	// TODO something doesn't work here
	return "unknown"
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
	sql := "INSERT INTO exif (photo_id, value_type, tag, string, num, denom, datetime, floating, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	if record.IsPersisted() {
		// TODO do we ever update exif tags?
		return record, fmt.Errorf("exif tags cannot be updated")
	}
	record.Timestamps = JustCreated(e.clock)

	err := e.db.QueryRowx(
		sql,
		photoID,
		record.Type,
		record.Tag,
		record.StringValue,
		record.Num,
		record.Denominator,
		record.DateTime,
		record.Floating,
		e.clock(),
		e.clock(),
	).Scan(&record.ID)

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
		record.StringValue = record.String()
		record.TypeName = typeNames[tiff.DataType(record.Type)]
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

var typeNames = map[tiff.DataType]string{
	tiff.DTByte:      "byte",
	tiff.DTAscii:     "ascii",
	tiff.DTShort:     "short",
	tiff.DTLong:      "long",
	tiff.DTRational:  "rational",
	tiff.DTSByte:     "signed byte",
	tiff.DTUndefined: "undefined",
	tiff.DTSShort:    "signed short",
	tiff.DTSLong:     "signed long",
	tiff.DTSRational: "signed rational",
	tiff.DTFloat:     "float",
	tiff.DTDouble:    "double",
}
