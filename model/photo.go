package model

import "github.com/ilikeorangutans/phts/db"

type Photo struct {
	db.PhotoRecord
	Renditions Renditions
	Exif       []ExifTag
	Collection Collection
}
