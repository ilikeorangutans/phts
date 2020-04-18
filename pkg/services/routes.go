package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilikeorangutans/phts/version"
	"github.com/ilikeorangutans/phts/web"
)

var Routes = []web.Section{
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
	},
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
