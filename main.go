package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

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
	}

	r := mux.NewRouter()
	web.BuildRoutes(r, phtsRoutes, []web.Filter{web.LoggingHandler, AddServicesToContext(db, backend)})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(bind, r)
}

func adminHomeHandler(w http.ResponseWriter, r *http.Request) {
	var admin = template.Must(template.ParseFiles("template/admin/base.tmpl", "template/admin/index.tmpl"))
	admin.Execute(w, nil)
}

func requireAdminAuth(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("requireAdminAuth()")
		wrap(w, r)
	}
}
