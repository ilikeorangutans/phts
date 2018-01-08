package public

import (
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
)

type CollectionRepository interface {
	FindByID(id int64) (model.Collection, error)
	FindBySlug(slug string) (model.Collection, error)
}

func NewCollectionRepository(dbx db.DB) CollectionRepository {
	return &publicCollectionRepo{
		collectionDB: db.NewCollectionDB(dbx),
	}
}

type publicCollectionRepo struct {
	collectionDB db.CollectionDB
}

func (r *publicCollectionRepo) FindByID(id int64) (model.Collection, error) {
	if record, err := r.collectionDB.FindByID(id); err != nil {
		return model.Collection{}, err
	} else {
		return model.Collection{
			CollectionRecord: record,
		}, nil
	}
}

func (r *publicCollectionRepo) FindBySlug(slug string) (model.Collection, error) {
	if record, err := r.collectionDB.FindBySlug(slug); err != nil {
		return model.Collection{}, err
	} else {
		return model.Collection{
			CollectionRecord: record,
		}, nil
	}
}
