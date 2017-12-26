package model

import "github.com/ilikeorangutans/phts/db"

type ShareSiteRepository interface {
	FindByID(int64) (ShareSite, error)
	List() ([]ShareSite, error)
	Save(ShareSite) (ShareSite, error)
	FindByDomain(domain string) (ShareSite, error)
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

func (r *shareSiteRepoImpl) FindByID(id int64) (ShareSite, error) {

	record, err := r.shareSiteDB.FindByID(id)

	result := ShareSite{
		ShareSiteRecord: record,
	}

	return result, err
}

func (r *shareSiteRepoImpl) FindByDomain(domain string) (ShareSite, error) {
	var result ShareSite
	record, err := r.shareSiteDB.FindByDomain(domain)
	if err != nil {
		return result, err
	}

	result.ShareSiteRecord = record

	return result, nil
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

func (r *shareSiteRepoImpl) Save(shareSite ShareSite) (ShareSite, error) {
	record, err := r.shareSiteDB.Save(shareSite.ShareSiteRecord)
	shareSite.ShareSiteRecord = record
	return shareSite, err
}
