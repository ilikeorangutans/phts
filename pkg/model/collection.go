package model

import "github.com/ilikeorangutans/phts/db"

// Collection is a single database level record of a collection.
type Collection struct {
	db.Record
	db.Timestamps
	db.Sluggable
	Name       string `db:"name" json:"name"`
	PhotoCount int    `db:"photo_count" json:"photoCount"`
}
