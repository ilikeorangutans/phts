package integration

import (
	"database/sql"
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
)

type TXAsDBWrapper struct {
	tx *sqlx.Tx
}

func (t *TXAsDBWrapper) Close() error {
	return nil
}

func (t *TXAsDBWrapper) QueryRowx(sql string, args ...interface{}) *sqlx.Row {
	return t.tx.QueryRowx(sql, args...)
}

func (t *TXAsDBWrapper) QueryRow(sql string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(sql, args...)
}

func (t *TXAsDBWrapper) Queryx(sql string, args ...interface{}) (*sqlx.Rows, error) {
	return t.tx.Queryx(sql, args...)
}

func (t *TXAsDBWrapper) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(sql, args...)
}

func (t *TXAsDBWrapper) Beginx() (db.TX, error) {
	fakeTransaction := &FakeTransaction{}
	return fakeTransaction, nil
}

type FakeTransaction struct {
	*sqlx.Tx
}

func (f *FakeTransaction) Rollback() error {
	log.Println("rolling back fake transaction")
	return nil
}

func (f *FakeTransaction) Commit() error {
	log.Println("comitting fake transaction")
	return nil
}
