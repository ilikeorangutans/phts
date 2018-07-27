package public

import (
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
)

type CollectionRepository interface {
	FindByID(id int64) (db.Collection, error)
	FindBySlug(slug string) (db.Collection, error)
}

func NewPublicCollectionRepository(dbx db.DB) model.CollectionFinder {
	return &publicCollectionRepo{
		collectionDB: db.NewCollectionDB(dbx),
	}
}

type publicCollectionRepo struct {
	collectionDB db.CollectionDB
}

func (r *publicCollectionRepo) FindByID(id int64) (db.Collection, error) {
	if record, err := r.collectionDB.FindByID(id); err != nil {
		return db.Collection{}, err
	} else {
		return record, nil
	}
}

func (r *publicCollectionRepo) FindBySlug(slug string) (db.Collection, error) {
	if record, err := r.collectionDB.FindBySlug(slug); err != nil {
		return db.Collection{}, err
	} else {
		return record, nil
	}
}
