package main

import (
	"net/http"

	"github.com/ilikeorangutans/phts/admin/collection"
	"github.com/ilikeorangutans/phts/admin/collection/photo"
	"github.com/ilikeorangutans/phts/web"
)

var phtsRoutes = []web.Section{
	{
		Path: "/admin",
		Middleware: []func(http.Handler) http.Handler{
			requireAdminAuth,
		},
		Routes: []web.Route{
			{
				Path:    "/",
				Handler: adminHomeHandler,
			},
			{
				Path:    "/collections",
				Handler: collection.IndexHandler,
			},
			{
				Path:    "/collections/new",
				Handler: collection.NewHandler,
			},
			{
				Path:    "/collections",
				Handler: collection.SaveHandler,
				Methods: []string{"POST"},
			},
			{
				Path:    "/collections/{slug:[a-z0-9-]+}",
				Handler: collection.ShowHandler,
				Middleware: []func(http.Handler) http.Handler{
					collection.RequireCollection,
				},
			},
			{
				Path:    "/collections/{slug:[a-z0-9-]+}/photos",
				Handler: collection.UploadPhotoHandler,
				Middleware: []func(http.Handler) http.Handler{
					collection.RequireCollection,
				},
				Methods: []string{"POST"},
			},
			{
				Path:    "/collections/{slug:[a-z0-9-]+}/photos/{photo_id:[0-9]+}",
				Handler: photo.ShowHandler,
				Middleware: []func(http.Handler) http.Handler{
					collection.RequireCollection,
					photo.RequirePhoto,
				},
			},
			{
				Path:    "/collections/{slug:[a-z0-9-]+}/photos/{photo_id:[0-9]+}",
				Handler: photo.DeleteHandler,
				Middleware: []func(http.Handler) http.Handler{
					collection.RequireCollection,
					photo.RequirePhoto,
				},
				Methods: []string{"POST", "DELETE"},
			},
			{
				Path:    "/collections/{slug:[a-z0-9-]+}/photos/renditions/{rendition_id:[0-9]+}",
				Handler: collection.ServeRendition,
				Middleware: []func(http.Handler) http.Handler{
					collection.RequireCollection,
				},
				Methods: []string{"GET"},
			},
		},
	},
	{
		Path: "/",
	},
}
