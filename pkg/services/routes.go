package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilikeorangutans/phts/version"
	"github.com/ilikeorangutans/phts/web"
)

func SetupServices(adminEmail, adminPassword string) []web.Section {
	return []web.Section{
		{
			Path: "/services",
			Routes: []web.Route{
				{
					Path:    "/ping",
					Handler: PingHandler,
				},
				{
					Path:    "/version",
					Handler: VersionHandler,
				},
			},
			Sections: []web.Section{
				{
					Path: "/internal",
					Middleware: []func(http.Handler) http.Handler{
						BasicAuth("services/internal", map[string]string{adminEmail: adminPassword}),
					},
					Routes: []web.Route{
						{
							Path:    "/",
							Handler: LandingPageHandler,
						},
					},
				},
			},
		},
	}
}

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	fmt.Fprintf(w, "pong")
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	version := struct {
		Sha       string `json:"sha"`
		BuildTime string `json:"buildTime"`
	}{
		Sha:       version.Sha,
		BuildTime: version.BuildTime,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(version)
}
