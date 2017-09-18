package model

import "github.com/ilikeorangutans/phts/db"

type RenditionRepository interface {
	FindByID(collection Collection, id int64) (Rendition, error)
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

func (r *renditionRepoImpl) FindByID(collection Collection, id int64) (rendition Rendition, err error) {
	record, err := r.renditionDB.FindByID(collection.ID, id)
	if err != nil {
		return rendition, err
	}

	return Rendition{
		RenditionRecord: record,
	}, nil
}
