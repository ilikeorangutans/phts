package model

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
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
		log.Printf("%v", t)
	}

	return result, nil
}

type ExifExtractor struct {
	tags []db.ExifRecord
}

func (extractor *ExifExtractor) Walk(name exif.FieldName, tag *tiff.Tag) error {
	exifTag, err := db.ExifRecordFromTiffTag(string(name), tag)
	log.Printf("%v", exifTag)
	if err != nil {
		log.Println(err)
	} else {
		extractor.tags = append(extractor.tags, exifTag)
	}
	return nil
}
