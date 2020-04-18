package main

import (
	"compress/gzip"
	"context"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/ilikeorangutans/phts/admin/api"
	"github.com/ilikeorangutans/phts/api/admin"
	"github.com/ilikeorangutans/phts/api/public"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/pkg/server"
	"github.com/ilikeorangutans/phts/pkg/services"
	"github.com/ilikeorangutans/phts/session"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/version"
	"github.com/ilikeorangutans/phts/web"
	"github.com/spf13/viper"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func AddServicesToContext(db db.DB, backend storage.Backend, sessions session.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "database", db)
			ctx = context.WithValue(ctx, "backend", backend)
			ctx = context.WithValue(ctx, "sessions", sessions)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func parseConfig() server.Config {
	return server.Config{
		Bind:             viper.GetString("bind"),
		DatabaseHost:     viper.GetString("db_host"),
		DatabaseUser:     viper.GetString("db_user"),
		DatabasePassword: viper.GetString("db_password"),
		DatabaseName:     viper.GetString("db_database"),
		DatabaseSSL:      viper.GetBool("db_ssl"),
		StorageEngine:    viper.GetString("storage_engine"),
		BucketName:       viper.GetString("minio_bucket"),
		MinioAccessKey:   viper.GetString("minio_access_key"),
		MinioSecretKey:   viper.GetString("minio_secret_key"),
		MinioEndpoint:    viper.GetString("minio_endpoint"),
		MinioUseSSL:      viper.GetBool("minio_use_ssl"),
	}
}

func setupEnvVars() {
	viper.SetEnvPrefix("phts")
	viper.AutomaticEnv()

	defaults := map[string]interface{}{
		"bind":        ":8080",
		"db_ssl":      false,
		"db_host":     "",
		"db_user":     "",
		"db_password": "",
		"db_database": "phts",

		"storage_engine": "file",

		"minio_bucket":     "",
		"minio_access_key": "",
		"minio_secret_key": "",
		"minio_endpoint":   "",
		"minio_use_ssl":    false,
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

func main() {
	log.Printf("phts starting up, version %s built on %s...", version.Sha, version.BuildTime)

	setupEnvVars()

	config := parseConfig()

	dbx, err := sqlx.Connect("postgres", config.DatabaseConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	ctx := context.Background()

	var backend storage.Backend
	switch config.StorageEngine {
	case "gcs":
		backend, err = storage.NewGCSBackend(config.ProjectID, ctx, config.BucketName)
	case "minio":
		backend, err = storage.NewMinIOBackend(config.MinioEndpoint, config.MinioAccessKey, config.MinioSecretKey, config.BucketName, config.MinioUseSSL)
	default:
		backend = storage.NewFileBackend("tmp")
	}
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(dbx.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrating database...")
	m, err := migrate.NewWithDatabaseInstance("file://db/migrate", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("Database up to date!")
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database migrated!")
	}

	exif.RegisterParsers(mknote.All...)

	wrappedDB := db.WrapDB(dbx)

	userDB := db.NewUserDB(wrappedDB)

	_, err = userDB.FindByEmail("admin@test.com")
	if err != nil {
		user := &db.UserRecord{
			Email: "admin@test.com",
		}

		err = user.UpdatePassword("test")
		if err != nil {
			panic(err)
		}

		err = userDB.Save(user)
		if err != nil {
			panic(err)
		}

		log.Printf("Created user %d %s", user.ID, user.Email)
	}

	sessionStorage := session.NewInMemoryStorage(30, time.Hour*1, time.Hour*24)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(AddServicesToContext(wrappedDB, backend, sessionStorage))
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"X-JWT", "Origin", "Accept", "Content-Type", "Cookie", "Content-Length", "Last-Modified", "Cache-Control"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "DELETE"},
		AllowCredentials: true,
		Debug:            false,
	})
	r.Use(cors.Handler)
	web.BuildRoutes(r, services.Routes, "/")
	web.BuildRoutes(r, adminAPIRoutes, "/")
	web.BuildRoutes(r, frontendAPIRoutes, "/")

	log.Println("Frontend files")
	setupFrontend(r, "/admin", "./ui-admin/dist")
	setupFrontend(r, "/", "./ui-public/dist")

	log.Printf("phts now waiting for requests on %s...", config.Bind)
	err = http.ListenAndServe(config.Bind, r)
	if err != nil {
		log.Fatal(err)
	}
}

func setupFrontend(r *chi.Mux, url string, dir string) {
	compression := middleware.Compress(gzip.DefaultCompression, "application/json", "application/javascript", "text/css")
	fileserver := http.FileServer(http.Dir(dir))

	// Serve anything under url/static from the directory. This is set during angular builds via -d.
	// This way we can easily distinguish between routing requests for angular (see below) and requests
	// for static assets.
	staticDir := filepath.Join(url, "static")
	staticDirPattern := filepath.Join(staticDir, "*")
	log.Printf("  GET %s", staticDirPattern)
	r.With(compression).Handle(staticDirPattern, http.StripPrefix(staticDir, fileserver))

	handlers := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{
			filepath.Join(url, "ngsw-worker.js"),
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filepath.Join(dir, "ngsw-worker.js"))
			},
		},
		{
			filepath.Join(url, "ngsw.json"),
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filepath.Join(dir, "ngsw.json"))
			},
		},
		{
			filepath.Join(url, "*"),
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filepath.Join(dir, "index.html"))
			},
		},
	}

	for _, location := range handlers {
		log.Printf("  GET %s", location.pattern)
		r.With(compression).HandleFunc(location.pattern, location.handler)
	}
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

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func panicHandler(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered %s: %s", r, debug.Stack())
			}
		}()
		wrap(w, req)
	}
}

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
