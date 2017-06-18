package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"html/template"

	"github.com/gorilla/mux"
	"github.com/ilikeorangutans/phts/web"
)

func SessionInViewHandler(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Opening session")
		ctx := context.WithValue(r.Context(), "", "")

		r = r.WithContext(ctx)

		wrap(w, r)

	}
}

func main() {
	bind := "localhost:8080"

	r := mux.NewRouter()
	web.BuildRoutes(r, phtsRoutes, []web.Filter{web.LoggingHandler, SessionInViewHandler})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	db, err := sql.Open("postgres", "user=jakob dbname=jakob sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrating database...")
	m, err := migrate.NewWithDatabaseInstance("file://db/migrate", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

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
