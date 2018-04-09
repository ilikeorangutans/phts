package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Queries interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRowx(string, ...interface{}) *sqlx.Row
	QueryRow(string, ...interface{}) *sql.Row
	Queryx(string, ...interface{}) (*sqlx.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Select(interface{}, string, ...interface{}) error
}

// DB is a type that represents a DB. Needed to mock out DBs in tests.
type DB interface {
	Queries
	Close() error
	Beginx() (TX, error)
	Rebind(string) string
}

// TX abstracts go's transaction type. Needed to mock out transactions in tests.
type TX interface {
	Queries
	Rollback() error
	Commit() error
}

// WrapDB wraps a given db instance into our own DB type.
func WrapDB(wrap *sqlx.DB) DB {
	return &DBWrapper{
		db: wrap,
	}
}

type DBWrapper struct {
	db *sqlx.DB
}

func (d *DBWrapper) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(sql, args...)
}

func (d *DBWrapper) QueryRowx(sql string, args ...interface{}) *sqlx.Row {
	return d.db.QueryRowx(sql, args...)
}

func (d *DBWrapper) Select(dest interface{}, sql string, args ...interface{}) error {
	return d.db.Select(dest, sql, args...)
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

func (d *DBWrapper) Rebind(s string) string {
	return d.db.Rebind(s)
}
