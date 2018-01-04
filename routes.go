package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/admin/api"
	"github.com/ilikeorangutans/phts/api/public"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/web"
)

var frontendAPIRoutes = []web.Section{
	{
		Path: "/api",
		Routes: []web.Route{
			{
				Path:    "/share/{slug:[A-Za-z0-9-]+}",
				Handler: public.ViewShareHandler,
			},
			{
				Path:    "/share/{slug:[A-Za-z0-9-]+}/renditions/{renditionID:[0-9]+}",
				Handler: public.ServeShareRenditionHandler,
				Methods: []string{"GET", "HEAD"},
			},
		},
		Middleware: []func(http.Handler) http.Handler{
			checkShareSite,
		},
	},
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
