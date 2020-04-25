package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/ilikeorangutans/phts/session"
)

const (
	ServicesInternalSessionCookieName = "PHTS_SERVICES_INTERNAL_SESSION_ID"
)

func RequiresAuthentication(sessions session.Storage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(ServicesInternalSessionCookieName)
			if err != nil {
				log.Printf("no cookie")
				http.Redirect(w, r, "/services/internal/login", http.StatusFound)
				return
			}

			sessionID := cookie.Value
			if !sessions.Check(sessionID) {
				log.Printf("no session %s", sessionID)
				http.Redirect(w, r, "/services/internal/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func LogoutHandler(sessions session.Storage, usersRepo *ServiceUsersRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func AuthenticationHandler(sessions session.Storage, usersRepo *ServiceUsersRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO clear any existing sessions
		// TODO disregard session ids coming from the client

		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			http.Error(w, "could not parse form", http.StatusBadRequest)
			return
		}

		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		user, err := usersRepo.FindByEmail(email)
		if err != nil {
			// TODO add error message to request
			log.Printf("error looking up user for authentication %s: %s", email, err)
			http.Redirect(w, r, "/services/internal/login", http.StatusFound)
			return
		}

		if !user.CheckPassword(password) {
			// TODO add error message to request
			log.Printf("wrong password")
			http.Redirect(w, r, "/services/internal/login", http.StatusFound)
			return
		}

		sessionID, err := security.GenerateRandomString(32)
		if err != nil {
			http.Error(w, "could not generate random string", http.StatusInternalServerError)
			return
		}

		// TODO sessions has an expiry but might be nice to explicitly set it here
		sessions.Add(sessionID, nil)

		_, err = usersRepo.JustLoggedIn(user)
		if err != nil {
			http.Error(w, "could not generate random string", http.StatusInternalServerError)
			return
		}

		// TODO set expiry date
		cookie := http.Cookie{
			Name:     ServicesInternalSessionCookieName,
			Value:    sessionID,
			Path:     "/services/internal",
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/services/internal/", http.StatusFound)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}
