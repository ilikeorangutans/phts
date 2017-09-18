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
	Path         string
	Handler      http.HandlerFunc
	Middleware   []func(http.Handler) http.Handler
	Methods      []string
	InSectionNav bool
}

func BuildRoutes(router chi.Router, sections []Section, base string) {
	for _, section := range sections {
		log.Printf("Section %s", filepath.Join(base, section.Path))
		subrouter := chi.NewRouter()
		router.Mount(section.Path, subrouter)
		subrouter.Use(section.Middleware...)

		for _, route := range section.Routes {

			methods := route.Methods
			if len(methods) == 0 {
				methods = []string{"GET"}
			}

			for _, m := range methods {
				switch m {
				case "GET":
					subrouter.With(route.Middleware...).Get(route.Path, route.Handler)
				case "POST":
					subrouter.With(route.Middleware...).Post(route.Path, route.Handler)
				case "HEAD":
					subrouter.With(route.Middleware...).Head(route.Path, route.Handler)
				default:
					log.Panicf("Don't know how to create route for method %s", m)
				}
			}

			fullPath := filepath.Join(base, section.Path, route.Path)
			log.Printf("  route %s %s", strings.Join(methods, ","), fullPath)
		}

		BuildRoutes(subrouter, section.Sections, filepath.Join(base, section.Path))

	}
}
