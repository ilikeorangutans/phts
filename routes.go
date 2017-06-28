package main

import (
	"github.com/ilikeorangutans/phts/admin/collection"
	"github.com/ilikeorangutans/phts/admin/collection/photo"
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
			{
				Path:    "/collections/{slug:[a-z0-9-]+}/photos/{photo_id:[0-9]+}",
				Handler: photo.ShowHandler,
				Filters: []web.Filter{
					photo.RequirePhoto,
					collection.RequireCollection,
				},
			},
			{
				Path:    "/collections/{slug:[a-z0-9-]+}/photos/{photo_id:[0-9]+}",
				Handler: photo.DeleteHandler,
				Filters: []web.Filter{
					photo.RequirePhoto,
					collection.RequireCollection,
				},
				Methods: []string{"POST", "DELETE"},
			},
			{
				Path:    "/collections/{slug:[a-z0-9-]+}/photos/renditions/{rendition_id:[0-9]+}",
				Handler: collection.ServeRendition,
				Filters: []web.Filter{
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
