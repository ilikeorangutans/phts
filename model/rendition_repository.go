package model

import (
	"fmt"

	"github.com/ilikeorangutans/phts/db"
)

type RenditionRepository interface {
	FindByID(collection Collection, id int64) (Rendition, error)
	FindByPhotoAndRenditionConfigurations(collection Collection, photo Photo, configs RenditionConfigurations) (Renditions, error)
}

func NewRenditionRepository(dbx db.DB) RenditionRepository {
	return &renditionRepoImpl{
		db:          dbx,
		renditionDB: db.NewRenditionDB(dbx),
	}
}

type renditionRepoImpl struct {
	db          db.DB
	renditionDB db.RenditionDB
}

func (r *renditionRepoImpl) FindByPhotoAndRenditionConfigurations(collection Collection, photo Photo, configs RenditionConfigurations) (renditions Renditions, err error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("no rendition configuration IDs provided")
	}
	var configIDs []int64
	for _, config := range configs {
		configIDs = append(configIDs, config.ID)
	}
	records, err := r.renditionDB.FindByPhotoAndConfigs(collection.ID, photo.ID, configIDs)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		renditions = append(renditions, Rendition{record, nil})
	}

	return renditions, nil
}

func (r *renditionRepoImpl) FindByID(collection Collection, id int64) (rendition Rendition, err error) {
	record, err := r.renditionDB.FindByID(collection.ID, id)
	if err != nil {
		return rendition, err
	}

	return Rendition{
		RenditionRecord: record,
	}, nil
}