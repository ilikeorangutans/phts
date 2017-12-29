package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/session"
)

func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessions := r.Context().Value("sessions").(session.Storage)
	values := make(map[string]interface{})
	values["id"] = user.ID
	values["date"] = time.Now().UTC().Unix()
	sessions.Add(tokenString, values)

	cookie := http.Cookie{
		Name:  "PHTS_ADMIN_JWT",
		Value: tokenString,
		/// Expires: time.Now().Add(time.Hour * 24), // TODO this should not be hardcoded
		Path: "/",
		//Domain:   strings.Split(r.Host, ":")[0],
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	resp := authenticationResponse{
		Email: user.Email,
		ID:    user.ID,
		JWT:   tokenString,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
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
