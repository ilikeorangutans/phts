package api

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/model"
)

func RequireCollection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		repo := model.CollectionRepoFromRequest(r)
		col, err := repo.FindBySlug(slug)
		log.Printf("RequireCollection %s", col)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "collection", col)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
