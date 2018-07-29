package model

import (
	"fmt"

	"github.com/ilikeorangutans/phts/db"
)

type RenditionConfigurationRepository interface {
	Save(*db.Collection, RenditionConfiguration) (RenditionConfiguration, error)
	Delete(RenditionConfiguration) error
	FindByIDs(ids []int64) (RenditionConfigurations, error)
}

func NewRenditionConfigurationRepository(dbx db.DB) RenditionConfigurationRepository {
	return &renditionConfigRepoImpl{
		db:                dbx,
		renditionConfigDB: db.NewRenditionConfigurationDB(dbx),
	}
}

type renditionConfigRepoImpl struct {
	db                db.DB
	renditionConfigDB db.RenditionConfigurationDB
}

func (r *renditionConfigRepoImpl) Save(collection *db.Collection, config RenditionConfiguration) (RenditionConfiguration, error) {
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

func (r *renditionConfigRepoImpl) FindByIDs(ids []int64) (RenditionConfigurations, error) {
	records, err := r.renditionConfigDB.FindByIDs(ids)
	if err != nil {
		return nil, err
	}

	var result RenditionConfigurations
	for _, record := range records {
		result = append(result, RenditionConfiguration{record})
	}

	return result, nil
}
