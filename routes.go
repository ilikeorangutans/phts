package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/admin/api"
	"github.com/ilikeorangutans/phts/api/admin"
	"github.com/ilikeorangutans/phts/api/public"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/pkg/services"
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
		Path: "/api/admin/authenticate",
		Routes: []web.Route{
			{
				Path:    "/",
				Handler: api.AuthenticateHandler,
				Methods: []string{"POST"},
			},
		},
	},
	{
		Path: "/api/admin",
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
				Path: "/account",
				Routes: []web.Route{
					{
						Path:    "/password",
						Handler: admin.UpdatePasswordHandler,
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
				},
				Sections: []web.Section{
					{
						Path: "/{slug:[a-z0-9-]+}",
						Middleware: []func(http.Handler) http.Handler{
							api.RequireCollection,
						},
						Routes: []web.Route{
							{
								Path:    "/",
								Handler: api.ShowCollectionHandler,
							},
							{
								Path:    "/",
								Handler: api.DeleteCollectionHandler,
								Methods: []string{"DELETE"},
							},
							{
								Path:    "/photos/recent",
								Handler: api.ListRecentPhotosHandler,
							},
							{
								Path:    "/photos/{id:[0-9]+}",
								Handler: api.ShowPhotoHandler,
							},
							{
								Path:    "/photos/{id:[0-9]+}",
								Handler: api.DeletePhotoHandler,
								Methods: []string{"DELETE"},
							},
							{
								Path:    "/photos/renditions/{id:[0-9]+}",
								Handler: api.ServeRenditionHandler,
								Methods: []string{"GET", "HEAD"},
							},
							{
								Path:    "/photos/{id:[0-9]+}/shares",
								Handler: api.ShowPhotoSharesHandler,
								Methods: []string{"GET"},
							},
							{
								Path:    "/photos/{id:[0-9]+}/shares",
								Handler: api.CreatePhotoShareHandler,
								Methods: []string{"POST"},
							},
							{
								Path:    "/photos",
								Handler: api.UploadPhotoHandler,
								Methods: []string{
									"POST",
								},
							},
							{
								Path:    "/photos",
								Handler: api.ListPhotosHandler,
							},
							{
								Path:    "/albums",
								Handler: api.ListAlbumsHandler,
							},
							{
								Path:    "/albums",
								Handler: api.CreateAlbumHandler,
								Methods: []string{"POST"},
							},

							{
								Path:    "/rendition_configurations",
								Handler: api.ListRenditionConfigurationsHandler,
							},
							{
								Path:    "/rendition_configurations",
								Handler: api.CreateRenditionConfigurationHandler,
								Methods: []string{"POST"},
							},
						},
						Sections: []web.Section{
							{
								Path: "/albums/{albumID:[0-9]+}",
								Middleware: []func(http.Handler) http.Handler{
									api.RequireAlbum,
								},
								Routes: []web.Route{
									{
										Path:    "/",
										Handler: api.AlbumDetailsHandler,
									},
									{
										Path:    "/",
										Handler: api.DeleteAlbumHandler,
										Methods: []string{"DELETE"},
									},
									{
										Path:    "/",
										Handler: api.UpdateAlbumHandler,
										Methods: []string{"POST"},
									},
									{
										Path:    "/photos",
										Handler: api.AlbumListPhotosHandler,
									},
									{
										Path:    "/photos",
										Handler: api.AddPhotosToAlbumHandler,
										Methods: []string{"POST"},
									},
								},
							},
						},
					},
				},
			},
		},
		Routes: []web.Route{
			{
				Path:    "/version",
				Handler: services.VersionHandler,
			},
		},
	},
}

func checkShareSite(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Checking share site %s", r.Host)
		db := model.DBFromRequest(r)
		shareSiteRepo := model.NewShareSiteRepository(db)
		shareSite, err := shareSiteRepo.FindByDomain(r.Host)
		if err != nil {
			// TODO need better handling here
			http.NotFound(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "shareSite", shareSite)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
