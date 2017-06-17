package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"html/template"

	"github.com/gorilla/mux"
	"github.com/ilikeorangutans/phts/model"
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

	colRepo := &model.DummyCollectionRepository{}
	col, _ := colRepo.FindByID(1)

	fmt.Println(col)

	db, err := sql.Open("postgres", "user=dev dbname=phts_dev sslmode=verify-full")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(db.Stats())

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
