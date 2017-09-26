package model

import "github.com/ilikeorangutans/phts/db"

type Photo struct {
	db.PhotoRecord
	Renditions Renditions `json:"renditions"`
	Exif       []ExifTag
	Collection Collection `json:"-"`
}
