package main

import (
	"net/http"

	"github.com/ilikeorangutans/phts/admin/api"
	"github.com/ilikeorangutans/phts/web"
)

var adminAPIRoutes = []web.Section{
	{
		Path: "/admin/api",
		Middleware: []func(http.Handler) http.Handler{
			requireAdminAuth,
		},
		Sections: []web.Section{
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
