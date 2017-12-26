package model

import "github.com/ilikeorangutans/phts/db"

type Share struct {
	db.ShareRecord
	ShareSite ShareSite
}
