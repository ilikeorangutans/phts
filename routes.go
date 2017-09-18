package main

import (
	"net/http"

	"github.com/ilikeorangutans/phts/admin/api"
	"github.com/ilikeorangutans/phts/web"
)

var adminAPIRoutes = []web.Section{
	{
		Path:       "/admin/api",
		Middleware: []func(http.Handler) http.Handler{},
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
						Path:    "/{slug:[a-z0-9]+}/photos",
						Handler: api.UploadPhotoHandler,
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
						Methods: []string{
							"POST",
						},
					},
				},
			},
		},
		Routes: []web.Route{},
	},
}

var phtsRoutes = []web.Section{
	{
		Path: "/admin",
		Middleware: []func(http.Handler) http.Handler{
			requireAdminAuth,
		},
		Sections: []web.Section{
			{
				Path:   "/collections",
				Routes: []web.Route{},
			},
		},
		Routes: []web.Route{
			{
				Path:    "/",
				Handler: adminHomeHandler,
			},
		},
	},
	{
		Path: "/",
	},
}
