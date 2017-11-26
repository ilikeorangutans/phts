package model

import "github.com/ilikeorangutans/phts/db"

type ShareSiteRepository interface {
	List() ([]ShareSite, error)
}

func NewShareSiteRepository(dbx db.DB) ShareSiteRepository {
	return &shareSiteRepoImpl{
		db:          dbx,
		shareSiteDB: db.NewShareSiteDB(dbx),
	}
}

type shareSiteRepoImpl struct {
	db          db.DB
	shareSiteDB db.ShareSiteDB
}

func (r *shareSiteRepoImpl) List() ([]ShareSite, error) {
	records, err := r.shareSiteDB.List()
	if err != nil {
		return nil, err
	}
	result := []ShareSite{}
	for _, record := range records {
		result = append(result, ShareSite{record})
	}

	return result, nil
}
