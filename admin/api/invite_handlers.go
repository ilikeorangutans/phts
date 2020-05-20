package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/web"
)

func GetInviteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	inviteID := chi.URLParam(r, "invite")
	db := web.DBFromRequest(r)
	userRepo := model.NewUserRepo(db)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	user, err := userRepo.ByInviteID(ctx, inviteID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Email: user.Email,
		Token: inviteID,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

type activateInviteRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func ActivateInviteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	inviteID := chi.URLParam(r, "invite")
	db := web.DBFromRequest(r)

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var activateRequest activateInviteRequest
	err := decoder.Decode(&activateRequest)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// TODO we should check the password here

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userRepo := model.NewUserRepo(db)
	_, err = userRepo.ActivateInvite(ctx, inviteID, activateRequest.Email, activateRequest.Name, activateRequest.Password)
	if err != nil {
		log.Printf("error %+v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
