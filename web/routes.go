package web

import (
	"log"
	"net/http"

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

			for _, m := range methods {
				switch m {
				case "GET":
					//subrouter.HandleFunc(route.Path, chain(route.Handler, routeFilters...))
					subrouter.With(route.Middleware...).Get(route.Path, route.Handler)
				}
			}

			//r := subrouter.HandleFunc(route.Path, chain(route.Handler, routeFilters...))
			//r.Methods(methods...)

			//fullPath := filepath.Join(section.Path, route.Path)
			//log.Printf("  route %s %s (%d filters)", strings.Join(methods, ","), fullPath, len(routeFilters))
		}
	}
}
