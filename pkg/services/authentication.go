package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/ilikeorangutans/phts/pkg/session"
	"github.com/pkg/errors"
)

const (
	ServicesInternalSessionCookieName = "PHTS_SERVICES_INTERNAL_SESSION_ID"
)

// RequiresAuthentication returns a middleware that requires authentication
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
		cookie, err := r.Cookie(ServicesInternalSessionCookieName)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		sessions.Remove(cookie.Value)
		// TODO this doesn't actually end sessions for some reason

		http.Redirect(w, r, "/services/internal/login", http.StatusFound)
	}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Errors    []string `json:"errors"`
	SessionID string   `json:"session_id,omitempty"`
}

func authenticate(usersRepo *ServiceUsersRepo, email, password string) (ServiceUser, string, error) {
	user, err := usersRepo.FindByEmail(email)
	if err != nil {
		// TODO add error message to request
		return user, "", errors.Wrap(err, "could not find user by email")
	}

	if !user.CheckPassword(password) {
		return user, "", errors.Wrap(err, "wrong password")
	}

	sessionID, err := security.GenerateRandomString(32)
	if err != nil {
		return user, "", errors.Wrap(err, "could not generate random string")
	}

	_, err = usersRepo.JustLoggedIn(user)
	if err != nil {
		return user, "", errors.Wrap(err, "could not record log in status")
	}

	return user, sessionID, nil
}

func AuthenticationHandler(sessions session.Storage, usersRepo *ServiceUsersRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO clear any existing sessions
		// TODO disregard session ids coming from the client
		defer r.Body.Close()
		w.Header().Set("Cache-Control", "no-store")
		var email, password string
		isJSONRequest := r.Header.Get("content-type") == "application/json"
		if isJSONRequest {
			w.Header().Set("content-type", "application/json")
			var authRequest authRequest
			var authResponse = authResponse{}
			encoder := json.NewEncoder(w)
			if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
				log.Printf("could not decode auth request json: %v", err)
				encoder.Encode(authResponse)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			email = authRequest.Username
			password = authRequest.Password
		} else {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "could not parse form", http.StatusBadRequest)
				return
			}

			email = r.PostFormValue("email")
			password = r.PostFormValue("password")
		}

		_, sessionID, err := authenticate(usersRepo, email, password)
		if err != nil {
			if isJSONRequest {
				w.WriteHeader(http.StatusUnauthorized)
				var authResponse = authResponse{
					Errors: []string{"authentication failed"},
				}
				encoder := json.NewEncoder(w)
				encoder.Encode(authResponse)
				return
			} else {
				log.Printf("error authenticating %s: %s", email, err)
				http.Redirect(w, r, "/services/internal/login", http.StatusFound)
				return
			}
		}

		// TODO sessions has an expiry but might be nice to explicitly set it here
		sessions.Add(sessionID, nil)

		// TODO set expiry date
		cookie := http.Cookie{
			Name:     ServicesInternalSessionCookieName,
			Value:    sessionID,
			Path:     "/services/internal",
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(w, &cookie)

		if isJSONRequest {
			w.WriteHeader(http.StatusCreated)
			var authResponse = authResponse{
				Errors:    []string{},
				SessionID: sessionID,
			}
			encoder := json.NewEncoder(w)
			encoder.Encode(authResponse)
		} else {
			http.Redirect(w, r, "/services/internal/", http.StatusFound)
		}
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := LoginPageTmpl().Execute(w, nil)
	if err != nil {
		log.Printf("%s", err)
	}
}
