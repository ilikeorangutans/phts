package web

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type Section struct {
	Path      string
	Filters   []Filter
	Routes    []Route
	Sections  []Section
	Templates []string
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Filters []Filter
	Methods []string
}

func BuildRoutes(router chi.Router, sections []Section, parentFilters []Filter) {
	for _, section := range sections {
		log.Printf("Section %s", section.Path)
		subrouter := chi.NewRouter()
		router.Mount(section.Path, subrouter)
		sectionFilters := append(section.Filters, parentFilters...)
		// TODO use chi's middleware mechanism here for filters

		for _, route := range section.Routes {

			methods := route.Methods
			if len(methods) == 0 {
				methods = []string{"GET"}
			}
			routeFilters := append(route.Filters, sectionFilters...)

			for _, m := range methods {
				switch m {
				case "GET":
					//subrouter.HandleFunc(route.Path, chain(route.Handler, routeFilters...))
					subrouter.Get(route.Path, chain(route.Handler, routeFilters...))
				}
			}

			//r := subrouter.HandleFunc(route.Path, chain(route.Handler, routeFilters...))
			//r.Methods(methods...)

			fullPath := filepath.Join(section.Path, route.Path)
			log.Printf("  route %s %s (%d filters)", strings.Join(methods, ","), fullPath, len(routeFilters))
		}
	}
}

func chain(h http.HandlerFunc, funcs ...Filter) http.HandlerFunc {
	result := h
	for _, f := range funcs {
		result = f(result)
	}
	return result
}

type Filter func(http.HandlerFunc) http.HandlerFunc

func LoggingHandler(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Begin %s %s", r.Method, r.RequestURI)
		wrap(w, r)
		log.Printf("Done  %s %s in %s", r.Method, r.RequestURI, time.Since(start))
	}
}
