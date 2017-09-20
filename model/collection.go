package model

import (
	"github.com/ilikeorangutans/phts/db"
)

// Collection is the highest level of organization in phts.
type Collection struct {
	db.CollectionRecord
}
