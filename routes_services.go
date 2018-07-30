package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/services/admin"
	"github.com/ilikeorangutans/phts/session"
	"github.com/ilikeorangutans/phts/web"
)

var serviceRoutes = []web.Section{
	{
		Path: "/services",
		Routes: []web.Route{
			{
				Path:    "/ping",
				Handler: servicesPingHandler,
			},
			{
				Path:    "/admin/authenticate",
				Handler: servicesAuthenticateHandler,
			},
		},
		Sections: []web.Section{
			{
				Path: "/admin",
			},
		},
	},
}

func servicesAuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	usernameAndPassword := authRequest{}
	err := decoder.Decode(&usernameAndPassword)
	if err != nil {
		log.Printf("failed to decode username and password json: %s", err)
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	dbx := model.DBFromRequest(r)
	adminService := admin.NewAdminService(dbx)
	admin, err := adminService.FindByEmailAndPassword(usernameAndPassword.Email, usernameAndPassword.Password)
	if err != nil {
		log.Printf("email %s not found: %s", usernameAndPassword.Email, err)
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": admin.Email,
		"id":    admin.ID,
	})
	// TODO secret is not a good choice
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessions := r.Context().Value("sessions").(session.Storage)
	values := make(map[string]interface{})
	values["user_id"] = user.ID
	values["date"] = time.Now().UTC().Unix()
	sessions.Add(tokenString, values)
	// TODO this should all be in a separate type
}

type authRequest struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

type authenticationResponse struct {
	Email string `json:"email"`
	ID    int64  `json:"id"`
	JWT   string `json:"jwt"`
}

func servicesPingHandler(w http.ResponseWriter, r *http.Request) {
	// Nothing to be done, we just need to respond with 200 OK
}
