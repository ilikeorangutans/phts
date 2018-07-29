package model

import (
	"github.com/ilikeorangutans/phts/db"
)

type Share struct {
	db.ShareRecord
	ShareSite               ShareSite                `json:"shareSite"`
	RenditionConfigurations []RenditionConfiguration `json:"renditionConfigurations"`
	Photos                  []Photo                  `json:"photos"`
	Collection              *db.Collection
}
