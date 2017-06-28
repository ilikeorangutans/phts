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

	"github.com/gorilla/mux"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/web"
)

func AddServicesToContext(db *sqlx.DB, backend storage.Backend) web.Filter {
	return func(wrap http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "database", db)
			ctx = context.WithValue(ctx, "backend", backend)
			r = r.WithContext(ctx)
			wrap(w, r)
		}
	}
}

func main() {
	log.Println("phts starting up...")
	bind := "localhost:8080"

	db, err := sqlx.Open("postgres", "user=jakob dbname=phts_dev sslmode=disable")
	//db, err := sqlx.Open("postgres", "user=dev dbname=phts_dev sslmode=disable")
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

	r := mux.NewRouter()
	web.BuildRoutes(r, phtsRoutes, []web.Filter{panicHandler, web.LoggingHandler, AddServicesToContext(db, backend)})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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

func requireAdminAuth(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrap(w, r)
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
