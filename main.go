package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
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

type phtsConfig struct {
	bind             string
	databaseHost     string
	databaseUser     string
	databasePassword string
	databaseName     string
	databaseSSL      bool
	storageEngine    string
	bucketName       string
	projectID        string
	minioAccessKey   string
	minioSecretKey   string
	minioEndpoint    string
	minioUseSSL      bool
}

func (c phtsConfig) DatabaseConnectionString() string {
	ssl := "enable"
	if !c.databaseSSL {
		ssl = "disable"
	}
	return fmt.Sprintf("user=%s host=%s password=%s dbname=%s sslmode=%s", c.databaseUser, c.databaseHost, c.databasePassword, c.databaseName, ssl)
}

func parseConfig() phtsConfig {
	return phtsConfig{
		bind:             viper.GetString("bind"),
		databaseHost:     viper.GetString("db_host"),
		databaseUser:     viper.GetString("db_user"),
		databasePassword: viper.GetString("db_password"),
		databaseName:     viper.GetString("db_database"),
		databaseSSL:      viper.GetBool("db_ssl"),
		storageEngine:    viper.GetString("storage_engine"),
		bucketName:       viper.GetString("minio_bucket"),
		minioAccessKey:   viper.GetString("minio_access_key"),
		minioSecretKey:   viper.GetString("minio_secret_key"),
		minioEndpoint:    viper.GetString("minio_endpoint"),
		minioUseSSL:      viper.GetBool("minio_use_ssl"),
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
	switch config.storageEngine {
	case "gcs":
		backend, err = storage.NewGCSBackend(config.projectID, ctx, config.bucketName)
	case "minio":
		backend, err = storage.NewMinIOBackend(config.minioEndpoint, config.minioAccessKey, config.minioSecretKey, config.bucketName, config.minioUseSSL)
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

	log.Printf("phts now waiting for requests on %s...", config.bind)
	err = http.ListenAndServe(config.bind, r)
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
