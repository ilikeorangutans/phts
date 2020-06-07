package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/model"
	model2 "github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/web"
)

func RequireCollection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		repo := model.CollectionRepoFromRequest(r)
		col, err := repo.FindBySlug(slug)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "collection", col)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// RequireCollectionBySlug checks the current url params for a slug and looks up the collection with that slug and stores it in the context.
// If no such collection exists for the current user, 404 is returned.
func RequireCollectionBySlug(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		user, err := web.UserFromRequest(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		dbx := web.DBFromRequest(r)
		repo, _ := model2.NewCollectionRepo(dbx)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		col, err := repo.FindBySlugAndUser(ctx, dbx, slug, user)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		ctx = context.WithValue(r.Context(), web.CollectionKey, col)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
