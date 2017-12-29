package main

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

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

func main() {
	log.Println("phts starting up...")
	bind := "localhost:8080"

	dbx, err := sqlx.Connect("postgres", "user=phts host=127.0.0.1 password=secret dbname=phts sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	backend := &storage.FileBackend{BaseDir: "tmp"}
	backend.Init()

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
		AllowedHeaders:   []string{"X-JWT", "Origin", "Accept", "Content-Type", "Cookie"},
		AllowCredentials: true,
		Debug:            true,
	})
	r.Use(cors.Handler)
	web.BuildRoutes(r, adminAPIRoutes, "/")
	web.BuildRoutes(r, frontendAPIRoutes, "/")

	r.Handle("/admin/frontend/*", http.StripPrefix("/admin/frontend/", http.FileServer(http.Dir("static"))))

	log.Println("phts now waiting for requests...")
	err = http.ListenAndServe(bind, r)
	if err != nil {
		log.Fatal(err)
	}
}

func requireAdminAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("requireAdminAuth should validate jwt tokens")

		jwt := r.Header.Get("X-JWT")
		log.Printf("jwt from headers: %s", jwt)

		if len(jwt) == 0 {
			cookie, err := r.Cookie("PHTS_ADMIN_JWT")
			if err != nil {
				log.Printf("error retrieving cookie: %s", err)
			} else {
				log.Printf("got cookie: %s", cookie)
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

		sess := sessions.Get(jwt)
		ctx := context.WithValue(r.Context(), "userID", sess["id"])
		log.Printf("User %d logged in", sess["id"])

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
