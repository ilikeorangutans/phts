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

	"github.com/namsral/flag"

	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/session"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/version"
	"github.com/ilikeorangutans/phts/web"
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
}

func (c phtsConfig) DatabaseConnectionString() string {
	ssl := "enable"
	if !c.databaseSSL {
		ssl = "disable"
	}
	return fmt.Sprintf("user=%s host=%s password=%s dbname=%s sslmode=%s", c.databaseUser, c.databaseHost, c.databasePassword, c.databaseName, ssl)
}

func parseConfig() phtsConfig {
	bindPtr := flag.String("bind", "localhost:8080", "hostname and port to bind to (BIND)")
	dbHostPtr := flag.String("db-host", "", "database host to connect to (DB_HOST)")
	dbUserPtr := flag.String("db-user", "", "database user (DB_USER)")
	dbNamePtr := flag.String("db-name", "", "database name (DB_NAME)")
	dbSSLPtr := flag.Bool("db-ssl", false, "connect to database over ssl (DB_SSL)")
	dbPasswordPtr := flag.String("db-password", "", "database password (DB_PASSWORD)")
	storageEnginePtr := flag.String("storage", "file", "storage engine (STORAGE)")
	storageBucketPtr := flag.String("storage-bucket", "file", "storage engine (STORAGE_BUCKET)")
	storageProjectIDPtr := flag.String("storage-project-id", "file", "storage engine (STORAGE_PROJECT_ID)")
	flag.Parse()
	return phtsConfig{
		bind:             *bindPtr,
		databaseHost:     *dbHostPtr,
		databaseUser:     *dbUserPtr,
		databasePassword: *dbPasswordPtr,
		databaseName:     *dbNamePtr,
		databaseSSL:      *dbSSLPtr,
		storageEngine:    *storageEnginePtr,
		bucketName:       *storageBucketPtr,
		projectID:        *storageProjectIDPtr,
	}
}

func main() {
	log.Printf("phts starting up, version %s built on %s...", version.Sha, version.BuildTime)

	config := parseConfig()

	dbx, err := sqlx.Connect("postgres", config.DatabaseConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	var backend storage.Backend
	if config.storageEngine == "gcs" {
		ctx := context.Background()
		backend, err = storage.NewGCSBackend(config.projectID, ctx, config.bucketName)
	} else {
		backend = storage.NewFileBackend("tmp")
	}

	err = db.ApplyMigrations(dbx.DB)
	if err != nil {
		log.Fatal(err)
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
	web.BuildRoutes(r, serviceRoutes, "/")
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
