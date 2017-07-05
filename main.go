package main

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"

	"html/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/web"
)

func AddServicesToContext(db *sqlx.DB, backend storage.Backend) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "database", db)
			ctx = context.WithValue(ctx, "backend", backend)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func main() {
	log.Println("phts starting up...")
	bind := "localhost:8080"

	//db, err := sqlx.Open("postgres", "user=jakob dbname=phts_dev sslmode=disable")
	db, err := sqlx.Open("postgres", "user=dev dbname=phts_dev sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	backend := &storage.FileBackend{BaseDir: "tmp"}
	backend.Init()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(AddServicesToContext(db, backend))
	web.BuildRoutes(r, phtsRoutes, "/")

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("phts now waiting for requests...")
	err = http.ListenAndServe(bind, r)
	if err != nil {
		log.Fatal(err)
	}
}

func adminHomeHandler(w http.ResponseWriter, r *http.Request) {
	var admin = template.Must(template.ParseFiles("template/admin/base.tmpl", "template/admin/index.tmpl"))
	err := admin.Execute(w, nil)
	if err != nil {
		log.Panic(err)
	}
}

func requireAdminAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
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
