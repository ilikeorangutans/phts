package server

import (
	"compress/gzip"
	"context"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/ilikeorangutans/phts/api/public"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	newmodel "github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/ilikeorangutans/phts/pkg/services"
	"github.com/ilikeorangutans/phts/pkg/session"
	"github.com/ilikeorangutans/phts/pkg/smtp"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/web"

	"database/sql"
	godb "database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err).Msg("could not connect to database")
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

	if err := m.EnsureUser(m.config.InitialUser, m.config.InitialUserPassword); err != nil {
		return errors.WithStack(err)
	}

	renditionUpdateRequestQueue := make(chan newmodel.RenditionUpdateRequest, 100)
	StartRenditionUpdateQueueHandler(ctx, m.db, m.backend, renditionUpdateRequestQueue, 2, 30*time.Minute)

	if err := m.SetupWebServer(ctx, renditionUpdateRequestQueue); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Main) SetupWebServer(ctx context.Context, renditionUpdateRequestQueue chan newmodel.RenditionUpdateRequest) error {
	sessionStorage := session.NewInMemoryStorage(30, time.Hour*1, time.Hour*24)
	email := smtp.NewEmailSender(m.config.SmtpHost, m.config.SmtpPort, m.config.SmtpUser, m.config.SmtpPassword, m.config.SmtpFrom)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(AddServicesToContext(m.db, m.backend, sessionStorage, renditionUpdateRequestQueue))
	cors := cors.New(cors.Options{
		// Add AllowOriginFunc to dynamically check origins
		AllowedOrigins:   []string{"*"}, // I'm pretty sure this defeats the entire purpose of CORS
		AllowedHeaders:   []string{"Authorization", "Origin", "Accept", "Content-Type", "Cookie", "Content-Length", "Last-Modified", "Cache-Control"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "DELETE"},
		AllowCredentials: true,
		MaxAge:           3600,
		Debug:            false,
	})
	r.Use(cors.Handler)

	compression := middleware.Compress(gzip.DefaultCompression, "application/json", "application/javascript", "text/css")
	fileserver := http.FileServer(http.Dir("templates/services/internal/static/"))
	r.With(compression).Handle("/services/internal/static/*", http.StripPrefix("/services/internal/static/", fileserver))
	r.Handle("/favicon.ico", http.FileServer(http.Dir("static")))
	log.Printf("  GET %s", "/services/internal/static/*")

	web.BuildRoutes(r, services.SetupServices(sessionStorage, m.db, email, m.config.AdminEmail, m.config.AdminPassword, m.config.ServerURL), "/")
	secret := m.config.JWTSecret

	if secret == "" {
		log.Warn().Msg("No JWT Secret set, generating a random token")
		randomSecret, err := security.GenerateRandomString(64)
		if err != nil {
			return errors.Wrap(err, "could not generate JWT secret")
		}
		secret = randomSecret
	}
	web.BuildRoutes(r, AdminAPIRoutes(secret), "/")
	web.BuildRoutes(r, frontendAPIRoutes, "/")

	log.Debug().Msg("Frontend Files")
	setupFrontend(r, "/admin", m.config.AdminStaticFilePath)
	setupFrontend(r, "/", m.config.FrontendStaticFilePath)

	log.Printf("phts now waiting for requests on %s...", m.config.Bind)
	err := http.ListenAndServe(m.config.Bind, r)
	if err != nil {
		return errors.Wrap(err, "could not start web server")
	}

	return nil
}

func (m *Main) EnsureUser(email, password string) error {
	userRepo := newmodel.NewUserRepo(m.db)
	user, err := userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "could not look up user")
	}
	userExists := !errors.Is(err, sql.ErrNoRows)
	hashedPassword, err := security.NewPassword(password)
	if err != nil {
		return errors.Wrap(err, "could not hash password")
	}

	if userExists {
		if user.Password.Matches(password) {
			log.Printf("initial user [%d] %s up to date", user.ID, user.Email)
			return nil
		}

		log.Printf("initial user [%d] %s exists but requires password update", user.ID, user.Email)
		user.Password = hashedPassword
		_, err = userRepo.Update(user)
		if err != nil {
			return errors.Wrap(err, "could not update user")
		}
		log.Printf("initial user [%d] %s password updated", user.ID, user.Email)
	} else {
		user := newmodel.User{
			Email:    email,
			Password: hashedPassword,
		}
		_, err = userRepo.Create(user)
		if err != nil {
			return errors.Wrap(err, "could not create new user")
		}
		log.Printf("initial user [%d] %s created", user.ID, user.Email)
	}
	return nil
}

func (m *Main) EnsureAdminUser(email, password string) error {
	usersRepo := services.NewServiceUsersRepo(m.db)
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
	log.Debug().Msg("migrating database")
	migrations, err := migrate.NewWithDatabaseInstance("file://db/migrate", "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "could not migrate database")
	}

	err = migrations.Up()
	if err == migrate.ErrNoChange {
		log.Debug().Msg("database schema up to date")
	} else if err != nil {
		return errors.Wrap(err, "could not migrate database")
	} else {
		log.Debug().Msg("database schema migrated")
	}

	return nil
}

func AddServicesToContext(dbx *sqlx.DB, backend storage.Backend, sessions session.Storage, queue chan newmodel.RenditionUpdateRequest) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			wrappedDB := db.WrapDB(dbx)
			ctx := context.WithValue(r.Context(), "database", wrappedDB)
			ctx = web.AddDBToContext(ctx, dbx)
			ctx = context.WithValue(ctx, "backend", backend)
			ctx = context.WithValue(ctx, "sessions", sessions)
			ctx = web.AddStorageBackendToContext(ctx, backend)
			ctx = web.AddRenditionUpdateRequestQueueToContext(ctx, queue)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func setupFrontend(r *chi.Mux, url string, dir string) {
	log.Printf("serving %s from %s", url, dir)
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

		ss, err := newmodel.FindShareSiteByDomain(r.Context(), web.DBFromRequest(r), r.Host)
		if err != nil {
			// TODO need better handling here
			http.NotFound(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "shareSite", shareSite)
		ctx = context.WithValue(ctx, web.ShareSiteKey, ss)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
