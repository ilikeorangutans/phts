package model

import (
	"fmt"
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
)

type ShareRepository interface {
	FindByShareSiteAndSlug(ShareSite, string) (Share, error)
	// FindByPhoto returns all shares that exist for the given photo.
	FindByPhoto(Photo, database.Paginator) ([]Share, error)
	Publish(Share) (Share, error)
}

func NewShareRepository(dbx db.DB) ShareRepository {
	return &shareRepoImpl{
		db:                     dbx,
		shareDB:                db.NewShareDB(dbx),
		shareRenditionConfigDB: db.NewShareRenditionConfigurationDB(dbx),
		renditionConfigRepo:    NewRenditionConfigurationRepository(dbx),
	}
}

type shareRepoImpl struct {
	db                     db.DB
	shareDB                db.ShareDB
	shareRenditionConfigDB db.ShareRenditionConfigurationDB
	renditionConfigRepo    RenditionConfigurationRepository
}

func (r *shareRepoImpl) FindByPhoto(photo Photo, paginator database.Paginator) ([]Share, error) {
	shares := []Share{}
	records, err := r.shareDB.FindByPhoto(photo.ID)
	if err != nil {
		return shares, err
	}

	shareSiteRepo := NewShareSiteRepository(r.db)

	for _, record := range records {
		// TODO this is super inefficient, we should either cache or do batch calls
		shareSite, err := shareSiteRepo.FindByID(record.ShareSiteID)
		if err != nil {
			return shares, err
		}
		var renditionConfigurations RenditionConfigurations
		if configRecords, err := r.shareRenditionConfigDB.FindByShare(record.ID); err == nil {
			for _, configRecord := range configRecords {
				renditionConfigurations = append(renditionConfigurations, RenditionConfiguration{configRecord.RenditionConfiguration})
			}
		} else {
			return shares, err
		}

		share := Share{
			ShareRecord:             record,
			ShareSite:               shareSite,
			RenditionConfigurations: renditionConfigurations,
			Photos:                  []Photo{},
		}
		shares = append(shares, share)
	}

	return shares, nil
}

func (r *shareRepoImpl) Publish(share Share) (Share, error) {
	if len(share.Photos) > 1 {
		return Share{}, fmt.Errorf("don't know how to share more than one item yet")
	}
	share.ShareRecord.PhotoID = share.Photos[0].ID
	share.ShareRecord.CollectionID = share.Collection.ID
	share.ShareRecord.ShareSiteID = share.ShareSite.ID

	log.Printf("Saving %v", share.ShareRecord)

	shareRecord, err := r.shareDB.Save(share.ShareRecord)
	if err != nil {
		return share, err
	}
	share.ShareRecord = shareRecord

	var shareRenditionConfigs []db.ShareRenditionConfigurationRecord
	for _, config := range share.RenditionConfigurations {
		shareRenditionConfigs = append(shareRenditionConfigs, db.ShareRenditionConfigurationRecord{
			ShareID:                  share.ID,
			RenditionConfigurationID: config.ID,
		})
	}

	_, err = r.shareRenditionConfigDB.SetForShare(share.ID, shareRenditionConfigs)

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
