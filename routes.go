package main

import (
	"github.com/ilikeorangutans/phts/admin/collection"
	"github.com/ilikeorangutans/phts/web"
)

var phtsRoutes = []web.Section{
	{
		Path: "/admin",
		Filters: []web.Filter{
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
				Filters: []web.Filter{
					collection.RequireCollection,
				},
			},
			{
				Path:    "/collections/{slug:[a-z0-9-]+}/photos",
				Handler: collection.UploadPhotoHandler,
				Filters: []web.Filter{
					collection.RequireCollection,
				},
				Methods: []string{"POST"},
			},
		},
	},
	{
		Path: "/",
	},
}