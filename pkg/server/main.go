package server

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
	"github.com/ilikeorangutans/phts/pkg/services"
	"github.com/ilikeorangutans/phts/pkg/smtp"
	"github.com/ilikeorangutans/phts/session"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/web"

	godb "database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func NewMain(ctx context.Context, config Config) (*Main, error) {
	storage, err := config.StorageBackend(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not create Main")
	}

	dbx, err := sqlx.ConnectContext(ctx, "postgres", config.DatabaseConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	return &Main{
		backend: storage,
		db:      dbx,
		config:  config,
	}, nil
}

// Main is the phts server application.
type Main struct {
	backend storage.Backend
	db      *sqlx.DB
	config  Config
}

func (m *Main) Run(ctx context.Context) error {
	exif.RegisterParsers(mknote.All...)

	if err := m.MigrateDatabase(); err != nil {
		return errors.WithStack(err)
	}

	if err := m.EnsureAdminUser(m.config.AdminEmail, m.config.AdminPassword); err != nil {
		return errors.WithStack(err)
	}

	if err := m.SetupWebServer(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Main) SetupWebServer() error {
	wrappedDB := db.WrapDB(m.db)
	sessionStorage := session.NewInMemoryStorage(30, time.Hour*1, time.Hour*24)
	email := smtp.NewEmailSender(m.config.SmtpHost, m.config.SmtpPort, m.config.SmtpUser, m.config.SmtpPassword, m.config.SmtpFrom)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(AddServicesToContext(wrappedDB, m.backend, sessionStorage))
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"X-JWT", "Origin", "Accept", "Content-Type", "Cookie", "Content-Length", "Last-Modified", "Cache-Control"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "DELETE"},
		AllowCredentials: true,
		Debug:            false,
	})
	r.Use(cors.Handler)
	web.BuildRoutes(r, services.SetupServices(sessionStorage, wrappedDB, email, m.config.AdminEmail, m.config.AdminPassword), "/")
	web.BuildRoutes(r, adminAPIRoutes, "/")
	web.BuildRoutes(r, frontendAPIRoutes, "/")

	log.Println("Frontend files")
	setupFrontend(r, "/admin", "./ui-admin/dist")
	setupFrontend(r, "/", "./ui-public/dist")

	log.Printf("phts now waiting for requests on %s...", m.config.Bind)
	err := http.ListenAndServe(m.config.Bind, r)
	if err != nil {
		return errors.Wrap(err, "could not start web server")
	}

	return nil
}

func (m *Main) EnsureAdminUser(email, password string) error {
	wrappedDB := db.WrapDB(m.db)

	usersRepo := services.NewServiceUsersRepo(wrappedDB)
	user, err := usersRepo.FindByEmail(email)
	if err == godb.ErrNoRows {
		user, err := usersRepo.NewUser(email, password, true)
		if err != nil {
			return errors.Wrap(err, "could not create admin user")
		}
		log.Printf("services/internal user [%d] %s created", user.ID, user.Email)
		return nil
	} else if err != nil {
		return errors.Wrap(err, "error while looking up admin user")
	}

	if user.CheckPassword(password) {
		log.Printf("services/internal user [%d] %s password up to date!", user.ID, email)
	} else {
		_, err := usersRepo.UpdatePassword(user, password)
		if err != nil {
			return errors.Wrap(err, "error ensuring admin user is up to date")
		}
		log.Printf("services/internal user [%d] %s password updated!", user.ID, email)
	}
	return nil
}

func (m *Main) MigrateDatabase() error {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "could not migrate database")
	}
	log.Println("Migrating database...")
	migrations, err := migrate.NewWithDatabaseInstance("file://db/migrate", "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "could not migrate database")
	}

	err = migrations.Up()
	if err == migrate.ErrNoChange {
		log.Println("Database up to date!")
	} else if err != nil {
		return errors.Wrap(err, "could not migrate database")
	} else {
		log.Println("Database migrated!")
	}

	return nil
}

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
