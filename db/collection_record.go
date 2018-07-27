package db

// Collection is a single database level record of a collection.
type Collection struct {
	Record
	Timestamps
	Sluggable
	Name       string `db:"name" json:"name"`
	PhotoCount int    `db:"photo_count" json:"photoCount"`
}
