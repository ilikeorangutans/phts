package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
)

func AuthenticateHandler(tokenForUser func(int64, string) (string, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		usernameAndPassword := authRequest{}
		err := decoder.Decode(&usernameAndPassword)
		if err != nil {
			log.Printf("failed to decode username and password json: %s", err)
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}

		log.Printf("authentication request for %s", usernameAndPassword.Username)

		dbx := model.DBFromRequest(r)
		userDB := db.NewUserDB(dbx)
		user, err := userDB.FindByEmail(usernameAndPassword.Username)
		if err != nil {
			log.Printf("username %s not found: %s", usernameAndPassword.Username, err)
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}

		log.Printf("user %d %s found", user.ID, user.Email)

		if !user.CheckPassword(usernameAndPassword.Password) {
			log.Printf("invalid password for user %s", user.Email)
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}

		log.Printf("user %d %s successfully authenticated", user.ID, user.Email)

		tokenString, err := tokenForUser(user.ID, user.Email)
		if err != nil {
			log.Printf("could not create jwt token: %+v", err)
			http.Error(w, "could not create JWT token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		resp := authenticationResponse{
			Email: user.Email,
			ID:    user.ID,
			JWT:   tokenString,
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(resp)
	}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authenticationResponse struct {
	Email string `json:"email"`
	ID    int64  `json:"id"`
	JWT   string `json:"jwt"`
}
