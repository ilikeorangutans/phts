package web

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/pkg/model"
	"github.com/jmoiron/sqlx"
)

const (
	DatabaseKey = iota
	WrappedDatabaseKey
	BackendKey
	SessionsKey
	UserKey
)

func AddUserToContext(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func UserFromRequest(r *http.Request) (model.User, error) {
	user := r.Context().Value(UserKey)
	if user == nil {
		return model.User{}, errors.New("no user in context")
	}

	return user.(model.User), nil
}

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
