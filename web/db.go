package web

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
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
		log.Fatal("DBFromRequest Could not get database from request, wrong type")
	}

	return db
}

func AddCollectionToContext(ctx context.Context, collection model.Collection) context.Context {
	return context.WithValue(ctx, CollectionKey, collection)
}

func CollectionFromRequest(r *http.Request) model.Collection {
	collection, ok := r.Context().Value(CollectionKey).(model.Collection)
	if !ok {
		log.Fatal("Could not get collection from request, wrong type")
	}

	return collection
}

func AddStorageBackendToContext(ctx context.Context, storage storage.Backend) context.Context {
	return context.WithValue(ctx, BackendKey, storage)
}

func StorageBackendFromRequest(r *http.Request) storage.Backend {
	storage, ok := r.Context().Value(BackendKey).(storage.Backend)
	if !ok {
		log.Fatal("Could not get storage backend from request, wrong type")
	}

	return storage
}

func AddRenditionUpdateRequestQueueToContext(ctx context.Context, queue chan model.RenditionUpdateRequest) context.Context {
	return context.WithValue(ctx, UpdateRenditionQueue, queue)
}

func GetRenditionUpdateRequestQueueFromRequest(r *http.Request) chan model.RenditionUpdateRequest {
	queue, ok := r.Context().Value(UpdateRenditionQueue).(chan model.RenditionUpdateRequest)
	if !ok {
		log.Fatal("no rendition update request queue in context")
	}
	return queue
}
