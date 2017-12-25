package db

type ShareRecord struct {
	Record
	Timestamps
	PhotoID     int64 `db:"photo_id"`
	ShareSiteID int64 `db:"share_site_id"`
}
