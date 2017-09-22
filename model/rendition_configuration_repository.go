package model

import (
	"fmt"

	"github.com/ilikeorangutans/phts/db"
)

type RenditionConfigurationRepository interface {
	Save(Collection, RenditionConfiguration) (RenditionConfiguration, error)
	Delete(RenditionConfiguration) error
}

func NewRenditionConfigurationRepository(dbx db.DB) RenditionConfigurationRepository {
	return &renditionConfigRepoImpl{
		db: dbx,
	}
}

type renditionConfigRepoImpl struct {
	db db.DB
}

func (r *renditionConfigRepoImpl) Save(collection Collection, config RenditionConfiguration) (RenditionConfiguration, error) {
	config.CollectionID = &collection.ID
	config.Private = false

	renditionConfigDB := db.NewRenditionConfigurationDB(r.db)

	var err error
	config.RenditionConfigurationRecord, err = renditionConfigDB.Save(config.RenditionConfigurationRecord)
	return config, err
}

func (r *renditionConfigRepoImpl) Delete(config RenditionConfiguration) error {
	if config.CollectionID == nil {
		return fmt.Errorf("cannot delete system configuration")
	}

	renditionConfigDB := db.NewRenditionConfigurationDB(r.db)
	return renditionConfigDB.Delete(config.ID)
}
