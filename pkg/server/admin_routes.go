package server

import (
	"context"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/admin/api"
	"github.com/ilikeorangutans/phts/api/admin"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	newmod "github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/pkg/services"
	"github.com/ilikeorangutans/phts/pkg/session"
	"github.com/ilikeorangutans/phts/web"
)

var adminAPIRoutes = []web.Section{
	{
		Path: "/api/admin/invite/{invite:[A-Za-z0-9-]+}",
		Routes: []web.Route{
			{
				Path:    "/",
				Methods: []string{"GET"},
				Handler: api.GetInviteHandler,
			},
			{
				Path:    "/",
				Methods: []string{"POST"},
				Handler: api.ActivateInviteHandler,
			},
		},
	},
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
				Path: "/photos",
				Routes: []web.Route{
					{
						Path:    "/",
						Handler: api.PhotoStreamHandler,
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
							api.RequireCollectionBySlug,
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

func requireAdminAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("requireAdminAuth should validate jwt tokens")

		jwt := r.Header.Get("X-JWT")

		if len(jwt) == 0 {
			cookie, err := r.Cookie("PHTS_ADMIN_JWT")
			if err != nil {
				log.Printf("error retrieving cookie: %s", err)
			} else {
				jwt = cookie.Value
			}
		}

		if len(jwt) == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO here we should validate the token

		sessions := r.Context().Value("sessions").(session.Storage)
		if !sessions.Check(jwt) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		db := r.Context().Value("database").(db.DB)
		sess := sessions.Get(jwt)

		userRepo := model.NewUserRepository(db)
		user, err := userRepo.FindByID(sess["user_id"].(int64))
		if err != nil {
			log.Printf("could not find user with id %v", sess["id"])
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", sess["id"])
		ctx = context.WithValue(r.Context(), "user", user)
		ctx = web.AddUserToContext(ctx, newmod.UserFromOldRecord(user))

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
