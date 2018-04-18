package db

import (
	_ "image/jpeg"
	_ "image/png"
)

type RenditionRecord struct {
	Record
	Timestamps

	PhotoID                  int64  `db:"photo_id" json:"photoID"`
	Original                 bool   `db:"original" json:"original"`
	Width                    uint   `db:"width" json:"width"`
	Height                   uint   `db:"height" json:"height"`
	Format                   string `db:"format" json:"format"`
	RenditionConfigurationID int64  `db:"rendition_configuration_id" json:"renditionConfigurationID"`
}
