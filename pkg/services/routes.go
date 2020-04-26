package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/session"
	"github.com/ilikeorangutans/phts/version"
	"github.com/ilikeorangutans/phts/web"
)

func SetupServices(sessions session.Storage, db db.DB, adminEmail, adminPassword string) []web.Section {
	serviceUsersRepo := &ServiceUsersRepo{
		db: db,
	}
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
				{
					Path:    "/internal/login",
					Handler: LoginHandler,
					Methods: []string{"POST", "GET"},
				},
				{
					Path:    "/internal/sessions/create",
					Handler: AuthenticationHandler(sessions, serviceUsersRepo),
					Methods: []string{"POST"},
				},
				// TODO add /internal/sessions/refresh and /internal/sessions/check
			},
			Sections: []web.Section{
				{
					Path: "/internal",
					Middleware: []func(http.Handler) http.Handler{
						RequiresAuthentication(sessions),
					},
					Routes: []web.Route{
						{
							Path:    "/",
							Handler: LandingPageHandler,
						},
						{
							Path:    "/sessions/destroy",
							Handler: LogoutHandler(sessions, serviceUsersRepo),
							Methods: []string{"GET", "POST", "DELETE"},
						},
						{
							Path:    "/service_users",
							Handler: ServiceUsersListHandler(serviceUsersRepo),
							Methods: []string{"GET"},
						},
					},
				},
			},
		},
	}
}

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := LandingPageTmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func ServiceUsersListHandler(usersRepo *ServiceUsersRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		paginator := ServiceUsersPaginator.PaginatorFromQuery(r.URL.Query())
		users, paginator, err := usersRepo.List(paginator)
		if err != nil {
			log.Printf("%+v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data := make(map[string]interface{})
		data["users"] = users
		data["paginator"] = paginator

		err = ServiceUsersPageTmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
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
