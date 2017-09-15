package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Queries interface {
	QueryRowx(string, ...interface{}) *sqlx.Row
	QueryRow(string, ...interface{}) *sql.Row
	Queryx(string, ...interface{}) (*sqlx.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

type DB interface {
	Queries
	Close() error
	Beginx() (TX, error)
}

type TX interface {
	Queries
	Rollback() error
	Commit() error
}

func WrapDB(wrap *sqlx.DB) DB {
	return &DBWrapper{
		db: wrap,
	}
}

type DBWrapper struct {
	db *sqlx.DB
}

func (d *DBWrapper) QueryRowx(sql string, args ...interface{}) *sqlx.Row {
	return d.db.QueryRowx(sql, args...)
}

func (d *DBWrapper) QueryRow(sql string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(sql, args...)
}

func (d *DBWrapper) Queryx(sql string, args ...interface{}) (*sqlx.Rows, error) {
	return d.db.Queryx(sql, args...)
}

func (d *DBWrapper) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(sql, args...)
}

func (d *DBWrapper) Close() error {
	return d.db.Close()
}

func (d *DBWrapper) Beginx() (TX, error) {
	return d.db.Beginx()
}
