package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/namsral/flag"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/session"
	"github.com/ilikeorangutans/phts/storage"
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
	log.Println("phts starting up...")

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
		user := db.UserRecord{
			Email: "admin@test.com",
		}

		err = user.UpdatePassword("test")
		if err != nil {
			panic(err)
		}

		user, err = userDB.Save(user)
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
	web.BuildRoutes(r, adminAPIRoutes, "/")
	web.BuildRoutes(r, frontendAPIRoutes, "/")

	js := http.FileServer(http.Dir("./ui/dist"))
	r.With(middleware.Compress(gzip.DefaultCompression, "application/json", "application/javascript", "text/css")).Handle("/static/*", http.StripPrefix("/static", js))
	r.HandleFunc("/ngsw-worker.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/dist/ngsw-worker.js")
	})
	r.With(middleware.Compress(gzip.DefaultCompression, "application/javascript")).HandleFunc("/ngsw.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/dist/ngsw.json")
	})
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/dist/index.html")
	})

	log.Printf("phts now waiting for requests on %s...", config.bind)
	err = http.ListenAndServe(config.bind, r)
	if err != nil {
		log.Fatal(err)
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
		user, err := userRepo.FindByID(sess["id"].(int64))
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
