package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ilikeorangutans/phts/admin/api"
	"github.com/ilikeorangutans/phts/api/admin"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/pkg/auth"
	newmod "github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/pkg/services"
	"github.com/ilikeorangutans/phts/web"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func AdminAPIRoutes(secret string) []web.Section {
	tokenForUser := func(id int64, email string) (string, error) {
		claim := auth.PhtsClaim{
			UserID:    id,
			UserEmail: email,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)

		signedToken, err := token.SignedString([]byte(secret))
		if err != nil {
			return "", errors.Wrap(err, "could not sign token")
		}

		return signedToken, nil
	}

	return []web.Section{
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
					Handler: api.AuthenticateHandler(tokenForUser),
					Methods: []string{"POST"},
				},
			},
		},
		{
			Path: "/api/admin",
			Middleware: []func(http.Handler) http.Handler{
				requireAdminAuthB(secret),
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
}

func requireAdminAuthB(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "no authorization header", http.StatusUnauthorized)
				return
			}

			token, err := jwt.ParseWithClaims(strings.TrimPrefix(authHeader, "Bearer "), &auth.PhtsClaim{}, func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil })
			if err != nil {
				log.Printf("invalid token %+v", err)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claim, ok := token.Claims.(*auth.PhtsClaim)
			if !ok {
				log.Printf("couldn't get claims")
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// TODO use new model
			db := r.Context().Value("database").(db.DB)
			userRepo := model.NewUserRepository(db)
			user, err := userRepo.FindByEmail(claim.UserEmail)
			if err != nil {
				log.Printf("user from claim not found")
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			ctx = web.AddUserToContext(ctx, newmod.UserFromOldRecord(user))

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
