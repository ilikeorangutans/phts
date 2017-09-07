package model

import (
	"fmt"

	"github.com/ilikeorangutans/phts/db"
)

type Collection struct {
	db.CollectionRecord
	collectionRepo CollectionRepository
}

func (c Collection) AddPhoto(filename string, data []byte) error {
	if !c.IsPersisted() {
		return fmt.Errorf("Cannot add photos to unpersisted collection")
	}

	return c.collectionRepo.AddPhoto(c, filename, data)
}
