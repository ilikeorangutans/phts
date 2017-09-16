package main

import (
	"net/http"

	"github.com/ilikeorangutans/phts/web"
)

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
