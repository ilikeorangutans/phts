package model

import "github.com/ilikeorangutans/phts/db"

type Photo struct {
	db.PhotoRecord
	Renditions Renditions `json:"renditions"`
	Exif       []ExifTag
	Collection Collection `json:"-"`
}

func NewPhotoFromRecord(record db.PhotoRecord, collection Collection, renditions Renditions) Photo {
	effectiveRenditions := renditions
	if effectiveRenditions == nil {
		effectiveRenditions = Renditions{}
	}
	return Photo{
		PhotoRecord: record,
		Renditions:  effectiveRenditions,
		Collection:  collection,
	}
}
