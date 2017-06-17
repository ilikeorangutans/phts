package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"html/template"

	"github.com/gorilla/mux"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/web"
)

func main() {
	bind := "localhost:8080"

	sections := []web.Section{
		{
			Path: "/admin",
			Filters: []web.Filter{
				requireAdminAuth,
			},
			Routes: []web.Route{
				{
					Path:    "/",
					Handler: adminHomeHandler,
				},
				{
					Path:    "/collections",
					Handler: collectionsIndexHandler,
					Filters: []web.Filter{
						requireCollectionFromSlug,
					},
				},
			},
		},
		{
			Path: "/",
		},
	}

	r := mux.NewRouter()
	web.BuildRoutes(r, sections, []web.Filter{loggingHandler})

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	colRepo := &model.DummyCollectionRepository{}
	col, _ := colRepo.FindByID(1)

	fmt.Println(col)

	http.ListenAndServe(bind, r)
}

func adminHomeHandler(w http.ResponseWriter, r *http.Request) {
	var admin = template.Must(template.ParseFiles("template/admin/base.tmpl", "template/admin/index.tmpl"))
	admin.Execute(w, nil)
}

func loggingHandler(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Begin %s %s", r.Method, r.RequestURI)
		wrap(w, r)
		log.Printf("Done  %s %s in %s", r.Method, r.RequestURI, time.Since(start))
	}
}

func requireAdminAuth(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("requireAdminAuth()")
		wrap(w, r)
	}
}

func requireCollectionFromSlug(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("requireCollectionFromSlug()")
		ctx := context.WithValue(r.Context(), "foo", "bar")
		r = r.WithContext(ctx)
		wrap(w, r)
	}
}

func collectionsIndexHandler(w http.ResponseWriter, r *http.Request) {
	foo, _ := r.Context().Value("foo").(string)
	log.Printf("Got foo: %s", foo)
	tmpl, err := template.ParseFiles("template/admin/base.tmpl", "template/admin/collection/index.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
	colRepo := &model.DummyCollectionRepository{}
	col, _ := colRepo.FindByID(1)

	data := make(map[string]interface{})
	data["collection"] = col
	tmpl.Execute(w, data)
}
