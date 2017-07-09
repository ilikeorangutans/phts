package web

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

type Section struct {
	Path       string
	Middleware []func(http.Handler) http.Handler
	Routes     []Route
	Sections   []Section
	Templates  []string
}

type Route struct {
	Path       string
	Handler    http.HandlerFunc
	Middleware []func(http.Handler) http.Handler
	Methods    []string
}

func BuildRoutes(router chi.Router, sections []Section) {
	for _, section := range sections {
		log.Printf("Section %s", section.Path)
		subrouter := chi.NewRouter()
		router.Mount(section.Path, subrouter)
		subrouter.Use(section.Middleware...)

		for _, route := range section.Routes {

			methods := route.Methods
			if len(methods) == 0 {
				methods = []string{"GET"}
			}
			//routeFilters := append(route.Filters, sectionFilters...)

			r := subrouter.With(route.Middleware...)
			for _, m := range methods {
				switch m {
				case "GET":
					r.Get(route.Path, route.Handler)
				case "POST":
					r.Post(route.Path, route.Handler)
				case "DELETE":
					r.Delete(route.Path, route.Handler)
				default:
					log.Panicf("Don't know %s", m)
				}
			}

			fullPath := filepath.Join(section.Path, route.Path)
			log.Printf("  route %s %s", strings.Join(methods, ","), fullPath)
		}
	}
}
