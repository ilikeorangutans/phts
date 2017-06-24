package model

import "github.com/ilikeorangutans/phts/db"

type Photo struct {
	db.PhotoRecord
	Renditions []Rendition
}

type Rendition struct {
	Original bool
	PhotoID  int64
	Data     []byte
}
