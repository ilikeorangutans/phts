package db

import "time"

type RenditionConfigurationRecord struct {
	Record
	Width        int       `db:"width"`
	Height       int       `db:"height"`
	Name         string    `db:"name"`
	Quality      int       `db:"quality"`
	CollectionID *int64    `db:"collection_id"`
	CreatedAt    time.Time `db:"created_at"`
}

type RenditionConfigurationDB interface {
	FindByID(id int64)
	Save(RenditionConfigurationRecord) (RenditionConfigurationRecord, error)
}
