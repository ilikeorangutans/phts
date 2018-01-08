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
	exifTimeLayout = "2006:01:02 15:04:05"
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
		log.Printf("error creating exif record from tag: %s", err.Error())
	} else {
		extractor.tags = append(extractor.tags, exifTag)
	}
	return nil
}

func ExifRecordFromTiffTag(name string, tag *tiff.Tag) (db.ExifRecord, error) {
	record := db.ExifRecord{
		Type: int(tag.Type),
		Tag:  strings.TrimRight(string(name), "\x00"),
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
			log.Printf("error getting string value from tag %s: %s", name, err.Error())
			return record, err
		} else {
			record.StringValue = strings.TrimRight(s, "\x00")
			// TODO sanitize input values
			weirdASCIIPrefix := "ASCII\x00\x00\x00"
			if strings.HasPrefix(record.StringValue, weirdASCIIPrefix) {
				record.StringValue = strings.TrimPrefix(record.StringValue, weirdASCIIPrefix)
			}

			if strings.Contains(name, "Date") {
				datetime, err := time.Parse(exifTimeLayout, record.StringValue)
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
	default:
		log.Printf("ignoring tag %s with unknown type %v", name, tag.Type)
	}

	return record, nil
}
