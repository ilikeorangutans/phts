package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/admin/api"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/web"
)

var frontendAPIRoutes = []web.Section{
	{
		Path: "/api",
		Routes: []web.Route{
			{
				Path:    "/share/{slug:[A-Za-z0-9-]+}",
				Handler: FrontendAPIShare,
			},
			{
				Path: "/share/{slug:[A-Za-z0-9-]+}/renditions/{id:[0-9]+}",
				// TODO need rendition serve handler here
				Handler: FrontendAPIShare,
			},
		},
		Middleware: []func(http.Handler) http.Handler{
			checkShareSite,
		},
	},
}

func FrontendIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("FrontendIndex"))
}

func FrontendShare(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("FrontendShare"))
}

func FrontendAPIShare(w http.ResponseWriter, r *http.Request) {
	log.Println("FrontendAPIShare")
	w.Header().Set("Content-Type", "application/json")

	db := model.DBFromRequest(r)
	storage := model.StorageFromRequest(r)
	shareSite := r.Context().Value("shareSite").(model.ShareSite)
	shareRepo := model.NewShareRepository(db)
	photoRepo := model.NewPhotoRepository(db, storage)
	collectionRepo := model.NewCollectionRepository(db, storage)

	type ShareResponse struct {
		Share  model.Share   `json:"share"`
		Photos []model.Photo `json:"photos"`
	}

	slug := chi.URLParam(r, "slug")

	share, err := shareRepo.FindByShareSiteAndSlug(shareSite, slug)
	if err != nil {
		log.Printf("No share found for slug %s and share site %s", slug, shareSite.Domain)
		http.NotFound(w, r)
		return
	}

	collection, err := collectionRepo.FindByID(share.CollectionID)
	if err != nil {
		log.Fatal(err)
	}
	photo, err := photoRepo.FindByID(collection, share.PhotoID)
	//if !photo.Published {
	//log.Printf("Photo %d not published", photo.ID)
	//http.NotFound(w, r)
	//return
	//}
	encoder := json.NewEncoder(w)

	// TODO we're dumping the entire photo record, need something smaller here
	resp := ShareResponse{
		Share:  share,
		Photos: []model.Photo{photo},
	}
	err = encoder.Encode(resp)
	if err != nil {
		log.Fatal(err)
	}
}

var adminAPIRoutes = []web.Section{
	{
		Path: "/admin/api/authenticate",
		Routes: []web.Route{
			{
				Path:    "/",
				Handler: api.AuthenticateHandler,
				Methods: []string{"POST"},
			},
		},
	},
	{
		Path: "/admin/api",
		Middleware: []func(http.Handler) http.Handler{
			requireAdminAuth,
		},
		Sections: []web.Section{
			{
				Path: "/share-sites",
				Routes: []web.Route{
					{
						Path:    "/",
						Handler: api.ListShareSitesHandler,
						Methods: []string{"GET"},
					},
					{
						Path:    "/",
						Handler: api.CreateShareSitesHandler,
						Methods: []string{"POST"},
					},
				},
			},
			{
				Path:       "/collections",
				Middleware: []func(http.Handler) http.Handler{},
				Routes: []web.Route{
					{
						Path:    "/",
						Handler: api.ListCollectionsHandler,
					},
					{
						Path:    "/",
						Handler: api.CreateCollectionHandler,
						Methods: []string{"POST"},
					},
					{
						Path:    "/{slug:[a-z0-9]+}",
						Handler: api.ShowCollectionHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/photos/recent",
						Handler: api.ListRecentPhotosHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/photos/{id:[0-9]+}",
						Handler: api.ShowPhotoHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/photos/renditions/{id:[0-9]+}",
						Handler: api.ServeRenditionHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
						Methods: []string{"GET", "HEAD"},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/photos/{id:[0-9]+}/shares",
						Handler: api.ShowPhotoSharesHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
						Methods: []string{"GET"},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/photos/{id:[0-9]+}/shares",
						Handler: api.CreatePhotoShareHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
						Methods: []string{"POST"},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/photos",
						Handler: api.UploadPhotoHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
						Methods: []string{
							"POST",
						},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/photos",
						Handler: api.ListPhotosHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/rendition_configurations",
						Handler: api.ListRenditionConfigurationsHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
					},
					{
						Path:    "/{slug:[a-z0-9]+}/rendition_configurations",
						Handler: api.CreateRenditionConfigurationHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
						Methods: []string{"POST"},
					},
				},
			},
		},
		Routes: []web.Route{},
	},
}

func checkShareSite(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Checking share site %s", r.Host)
		db := model.DBFromRequest(r)
		shareSiteRepo := model.NewShareSiteRepository(db)
		shareSite, err := shareSiteRepo.FindByDomain(r.Host)
		if err != nil {
			log.Printf("%s", err)
			// TODO need better handling here
			http.NotFound(w, r)
			return
		}
		log.Printf("Found share site %s", shareSite)
		ctx := context.WithValue(r.Context(), "shareSite", shareSite)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
