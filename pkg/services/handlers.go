package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/pkg/smtp"
	"github.com/ilikeorangutans/phts/version"
	"github.com/jordan-wright/email"
)

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := LandingPageTmpl().Execute(w, nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func SmtpTestHandler(emailer *smtp.Email, serverURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			recipient := r.PostFormValue("email")

			log.Printf("sending test email to %s", recipient)
			e := email.NewEmail()
			e.To = []string{recipient}
			e.Subject = "Test email from phts"
			e.Text = []byte(fmt.Sprintf("this is a test email from phts at %s", serverURL))
			err := emailer.Send(e)
			if err != nil {
				log.Printf("%+v", err)
			}
		}
		w.Header().Set("Content-Type", "text/html")
		data := make(map[string]interface{})
		data["settings"] = emailer
		err := SmtpTestTmpl().Execute(w, data)
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
		UserInviteEmailTmpl().Execute(&b, data)
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

		err = UsersPageTmpl().Execute(w, data)
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

		err = ServiceUsersPageTmpl().Execute(w, data)
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
