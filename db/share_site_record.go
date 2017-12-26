package db

import (
	"time"
)

type ShareSiteRecord struct {
	Record
	Timestamps
	Domain string `db:"domain" json:"domain"`
}

type ShareSiteDB interface {
	FindByID(int64) (ShareSiteRecord, error)
	FindByDomain(domain string) (ShareSiteRecord, error)
	Save(record ShareSiteRecord) (ShareSiteRecord, error)
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

func (s *shareSiteSQLDB) FindByID(id int64) (record ShareSiteRecord, err error) {
	sql := "SELECT * FROM share_sites WHERE id = $1 LIMIT 1"
	err = s.db.QueryRowx(
		sql,
		id,
	).StructScan(&record)

	return record, err
}

func (s *shareSiteSQLDB) FindByDomain(domain string) (record ShareSiteRecord, err error) {
	sql := "SELECT * FROM share_sites WHERE domain = $1 LIMIT 1"

	err = s.db.QueryRowx(
		sql,
		domain,
	).StructScan(&record)

	return record, err
}

func (s *shareSiteSQLDB) Save(record ShareSiteRecord) (ShareSiteRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(s.clock)
		sql := "UPDATE share_sites SET domain=$1, updated_at=$2 WHERE id=$3"
		err = checkResult(s.db.Exec(
			sql,
			record.Domain,
			record.UpdatedAt,
			record.ID,
		))
	} else {
		record.Timestamps = JustCreated(s.clock)
		sql := "INSERT INTO share_sites (domain, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id"

		err = s.db.QueryRow(
			sql,
			record.Domain,
			record.CreatedAt,
			record.UpdatedAt,
		).Scan(&record.ID)
	}
	return record, err
}
