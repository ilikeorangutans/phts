package db

import (
	"time"
)

type ShareSiteRecord struct {
	Record
	Timestamps

	Domain       string `db:"domain"`
	CollectionID int64  `db:"collection_id" json:"collectionID"`
}

type ShareSiteDB interface {
	FindByDomain(collectionID int64, domain string) (ShareSiteRecord, error)
	Save(collectionID int64, record ShareSiteRecord) (ShareSiteRecord, error)
	List() ([]ShareSiteRecord, error)
}

func NewShareSiteDB(db DB) ShareSiteDB {
	return &shareSiteSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type shareSiteSQLDB struct {
	db    DB
	clock Clock
}

func (s *shareSiteSQLDB) List() (records []ShareSiteRecord, err error) {
	sql := "SELECT * FROM share_sites"

	rows, err := s.db.Queryx(sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		record := ShareSiteRecord{}
		err = rows.StructScan(&record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, err
}

func (s *shareSiteSQLDB) FindByDomain(collectionID int64, domain string) (record ShareSiteRecord, err error) {

	sql := "SELECT * FROM share_sites WHERE collection_id = $1 AND domain = $2 LIMIT 1"

	err = s.db.QueryRowx(
		sql,
		collectionID,
		domain,
	).StructScan(&record)

	return record, err
}

func (s *shareSiteSQLDB) Save(collectionID int64, record ShareSiteRecord) (ShareSiteRecord, error) {
	record.CollectionID = collectionID
	var err error
	if record.IsPersisted() {
		record.JustUpdated(s.clock)
		sql := "UPDATE share_sites SET domain=$1, updated_at=$2 WHERE id=$3 AND collection_id=$4"
		err = checkResult(s.db.Exec(
			sql,
			record.Domain,
			record.UpdatedAt,
			record.ID,
			record.CollectionID,
		))
	} else {
		record.Timestamps = JustCreated(s.clock)
		sql := "INSERT INTO share_sites (domain, collection_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"

		err = s.db.QueryRow(
			sql,
			record.Domain,
			record.CollectionID,
			record.CreatedAt,
			record.UpdatedAt,
		).Scan(&record.ID)
	}
	return record, err
}
