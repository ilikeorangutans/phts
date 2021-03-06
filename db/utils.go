package db

import (
	sq "github.com/Masterminds/squirrel"
)

func queryAndStructScan(db DB, s sq.SelectBuilder, record interface{}) error {
	sql, args, err := s.ToSql()
	if err != nil {
		return err
	}
	return db.QueryRowx(sql, args...).StructScan(record)
}
