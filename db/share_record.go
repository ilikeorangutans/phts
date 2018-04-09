package db

import "fmt"

type ShareRecord struct {
	Record
	Timestamps
	PhotoID      int64  `db:"photo_id" json:"photoID"`
	CollectionID int64  `db:"collection_id" json:"collectionID"`
	ShareSiteID  int64  `db:"share_site_id" json:"shareSiteID"`
	Slug         string `db:"slug" json:"slug"`
}

type ShareRenditionConfigurationRecord struct {
	Timestamps
	ShareID                  int64                        `db:"share_id"`
	RenditionConfigurationID int64                        `db:"rendition_configuration_id"`
	RenditionConfiguration   RenditionConfigurationRecord `db:"rc"`
}

func (s ShareRenditionConfigurationRecord) IsPersisted() bool {
	return s.ShareID != 0 && s.RenditionConfigurationID != 0
}

func (s ShareRenditionConfigurationRecord) String() string {
	return fmt.Sprintf("ShareRenditionConfiguration{ShareID:%d,RenditionConfigurationID:%d}", s.ShareID, s.RenditionConfigurationID)
}
