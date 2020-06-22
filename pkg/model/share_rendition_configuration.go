package model

import "github.com/ilikeorangutans/phts/db"

type ShareRenditionConfigurationRecord struct {
	db.Timestamps
	ShareID                  int64 `db:"share_id"`
	RenditionConfigurationID int64 `db:"rendition_configuration_id"`
}
