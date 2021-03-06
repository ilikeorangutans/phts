package model

import (
	"time"

	"github.com/ilikeorangutans/phts/db"
)

type Photo struct {
	db.Record
	db.Timestamps

	CollectionID   int64      `db:"collection_id" json:"collectionID"`
	RenditionCount int        `db:"rendition_count" json:"renditionCount"`
	Description    string     `db:"description" json:"description"`
	Filename       string     `db:"filename" json:"filename"`
	TakenAt        *time.Time `db:"taken_at" json:"takenAt"`
	Published      bool       `db:"published" json:"published"`
}
