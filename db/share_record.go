package db

import "time"

type ShareRecord struct {
	Record
	Timestamps
	PhotoID      int64  `db:"photo_id" json:"photo_id"`
	CollectionID int64  `db:"collection_id" json:"collection_id"`
	ShareSiteID  int64  `db:"share_site_id" json:"share_site_id"`
	Slug         string `db:"slug" json:"slug"`
}

type ShareDB interface {
	FindByShareSiteAndSlug(shareSiteID int64, slug string) (ShareRecord, error)
	FindByPhoto(photoID int64) ([]ShareRecord, error)
	Save(ShareRecord) (ShareRecord, error)
}

func NewShareDB(dbx DB) ShareDB {
	return &shareSQLDB{
		db:    dbx,
		clock: time.Now,
	}
}

type shareSQLDB struct {
	db    DB
	clock Clock
}

func (c *shareSQLDB) Save(record ShareRecord) (ShareRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(c.clock)
		// TODO implement me
		sql := "XXX implement me"
		err = checkResult(c.db.Exec(
			sql,
			record.PhotoID,
			record.ShareSiteID,
			record.Slug,
			record.UpdatedAt,
			record.ID,
		))
	} else {
		record.Timestamps = JustCreated(c.clock)
		sql := "INSERT INTO shares (photo_id, share_site_id, slug, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"

		err = c.db.QueryRow(
			sql,
			record.PhotoID,
			record.ShareSiteID,
			record.Slug,
			record.CreatedAt,
			record.UpdatedAt,
		).Scan(&record.ID)
	}

	return record, err
}

func (c *shareSQLDB) FindByShareSiteAndSlug(shareSiteID int64, slug string) (ShareRecord, error) {
	sql := "SELECT * FROM shares WHERE share_site_id = $1 AND slug = $2 LIMIT 1"

	var record ShareRecord
	err := c.db.QueryRowx(sql, shareSiteID, slug).StructScan(&record)
	return record, err
}

func (c *shareSQLDB) FindByPhoto(photoID int64) ([]ShareRecord, error) {
	sql := "SELECT * FROM shares WHERE photo_id = $1"

	rows, err := c.db.Queryx(sql, photoID)

	var result []ShareRecord
	for rows.Next() {
		record := ShareRecord{}
		err = rows.StructScan(&record)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, err
}
