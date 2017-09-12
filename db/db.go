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
