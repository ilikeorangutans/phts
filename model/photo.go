package model

import (
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/metadata"
)

type Photo struct {
	db.PhotoRecord
	Renditions Renditions `json:"renditions"`
	Exif       []metadata.ExifTag
	//Collection db.Collection `json:"-"`
}

func NewPhotoFromRecord(record db.PhotoRecord, collection *db.Collection, renditions Renditions) Photo {
	effectiveRenditions := renditions
	if effectiveRenditions == nil {
		effectiveRenditions = Renditions{}
	}
	return Photo{
		PhotoRecord: record,
		Renditions:  effectiveRenditions,
		//Collection:  collection,
	}
}
