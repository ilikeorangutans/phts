package model

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

const (
	exifTimeLayout = "2006:01:02 15:04:0"
)

type ExifTags []ExifTag

func (e ExifTags) ByName(name string) (ExifTag, error) {
	for _, t := range e {
		if t.Tag == name {
			return t, nil
		}
	}

	return ExifTag{}, errors.New("No such tag")
}

type ExifTag struct {
	db.ExifRecord
}

func (e ExifTag) String() string {
	switch e.Type {
	case tiff.DTAscii:
		return fmt.Sprintf("%s", e.StringValue)
	case tiff.DTLong, tiff.DTShort:
		return fmt.Sprintf("%d", e.Num)
	case tiff.DTRational, tiff.DTSRational:
		return fmt.Sprintf("%d/%d", e.Num, e.Denominator)
	}
	return "unknown"
}

func ExifTagsFromPhoto(data []byte) (ExifTags, error) {
	x, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		log.Println("Decoding failed", err)
		return nil, err
	}

	extractor := &ExifExtractor{}
	x.Walk(extractor)

	result := []ExifTag{}
	for _, t := range extractor.tags {
		result = append(result, ExifTag{t})
	}

	return result, nil
}

type ExifExtractor struct {
	tags []db.ExifRecord
}

func (extractor *ExifExtractor) Walk(name exif.FieldName, tag *tiff.Tag) error {
	exifTag, err := ExifRecordFromTiffTag(string(name), tag)
	if err != nil {
		log.Println(err)
	} else {
		extractor.tags = append(extractor.tags, exifTag)
	}
	return nil
}

func ExifRecordFromTiffTag(name string, tag *tiff.Tag) (db.ExifRecord, error) {
	record := db.ExifRecord{
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

			if strings.Contains(name, "Date") {
				datetime, err := time.Parse(exifTimeLayout, record.StringValue)
				log.Printf("Parsed datetime: %s => %v, %v\n", name, datetime, err)
				if err == nil {
					record.DateTime = &datetime
					return record, nil
				}
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
