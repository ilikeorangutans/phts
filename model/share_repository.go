package model

import "github.com/ilikeorangutans/phts/db"

type ShareRepository interface {
	FindByShareSiteAndSlug(ShareSite, string) (Share, error)
	// TODO need paginator here
	FindByPhoto(Photo) ([]Share, error)
	Save(Share) (Share, error)
}

func NewShareRepository(dbx db.DB) ShareRepository {
	return &shareRepoImpl{
		db:      dbx,
		shareDB: db.NewShareDB(dbx),
	}
}

type shareRepoImpl struct {
	db      db.DB
	shareDB db.ShareDB
}

func (r *shareRepoImpl) FindByPhoto(photo Photo) ([]Share, error) {
	var shares []Share
	records, err := r.shareDB.FindByPhoto(photo.ID)
	if err != nil {
		return shares, err
	}

	shareSiteRepo := NewShareSiteRepository(r.db)

	for _, record := range records {
		// TODO this is super inefficient
		shareSite, err := shareSiteRepo.FindByID(record.ShareSiteID)
		if err != nil {
			return shares, err
		}
		share := Share{
			ShareRecord: record,
			ShareSite:   shareSite,
		}
		shares = append(shares, share)
	}

	return shares, nil
}

func (r *shareRepoImpl) Save(share Share) (Share, error) {
	record, err := r.shareDB.Save(share.ShareRecord)
	share.ShareRecord = record
	return share, err
}

func (r *shareRepoImpl) FindByShareSiteAndSlug(shareSite ShareSite, slug string) (Share, error) {
	var share Share

	record, err := r.shareDB.FindByShareSiteAndSlug(shareSite.ID, slug)
	if err != nil {
		return share, err
	}
	share.ShareRecord = record

	return share, nil
}
