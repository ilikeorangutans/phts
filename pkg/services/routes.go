package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/pkg/smtp"
	"github.com/ilikeorangutans/phts/session"
	"github.com/ilikeorangutans/phts/version"
	"github.com/ilikeorangutans/phts/web"
	"github.com/jordan-wright/email"
)

func SetupServices(sessions session.Storage, db db.DB, emailer *smtp.Email, adminEmail, adminPassword, serverURL string) []web.Section {
	serviceUsersRepo := NewServiceUsersRepo(db)
	usersRepo := model.NewUserRepo(db)

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
						},
						{
							Path:    "/users",
							Handler: UsersListHandler(usersRepo),
						},
						{
							Path:    "/users/invite",
							Handler: UsersInviteHandler(usersRepo, emailer, serverURL),
							Methods: []string{"POST"},
						},
						{
							Path:    "/smtp_test",
							Handler: SmtpTestHandler(emailer),
							Methods: []string{"GET", "POST"},
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

func SmtpTestHandler(emailer *smtp.Email) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			recipient := r.PostFormValue("email")

			log.Printf("sending test email to %s", recipient)
			e := email.NewEmail()
			e.To = []string{recipient}
			e.Subject = "Test email from phts"
			e.Text = []byte("this is a test email from phts")
			err := emailer.Send(e)
			if err != nil {
				log.Printf("%+v", err)
			}
		}
		w.Header().Set("Content-Type", "text/html")
		data := make(map[string]interface{})
		data["settings"] = emailer
		err := SmtpTestTmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

	}
}

func UsersInviteHandler(usersRepo *model.UserRepo, emailer *smtp.Email, serverURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		recipient := r.PostFormValue("email")
		log.Printf("inviting %s", recipient)

		// TODO if the user already exists, generate new token and resend.

		user, err := usersRepo.NewUser(recipient)
		if err != nil {
			log.Printf("%+v", err)
		}
		log.Printf("Created %v", user)

		e := email.NewEmail()
		e.To = []string{recipient}
		e.Subject = "You've been invited to phts"

		var b bytes.Buffer
		data := make(map[string]interface{})
		data["token"] = user.PasswordChangeToken
		data["email"] = user.Email
		data["server_url"] = serverURL
		UserInviteEmailTmpl.Execute(&b, data)
		e.Text = b.Bytes()
		err = emailer.Send(e)
		if err != nil {
			log.Printf("%+v", err)
		}

		http.Redirect(w, r, "/services/internal/users", http.StatusFound)
	}
}

func UsersListHandler(usersRepo *model.UserRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		paginator := ServiceUsersPaginator.PaginatorFromQuery(r.URL.Query())
		users, paginator, err := usersRepo.List(paginator)
		if err != nil {
			// TODO surface this error to the user
			log.Printf("%+v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data := make(map[string]interface{})
		data["users"] = users
		data["paginator"] = paginator

		err = UsersPageTmpl.Execute(w, data)
		if err != nil {
			log.Printf("%+v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

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
