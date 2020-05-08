package web

import (
	"context"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

const (
	DatabaseKey = iota
	WrappedDatabaseKey
	BackendKey
	SessionsKey
)

func AddDBToContext(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, DatabaseKey, db)
}

func DBFromRequest(r *http.Request) *sqlx.DB {
	db, ok := r.Context().Value(DatabaseKey).(*sqlx.DB)
	if !ok {
		log.Fatal("Could not get database from request, wrong type")
	}

	return db
}
