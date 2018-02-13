package admin

import "net/http"

type passwordChangeRequest struct {
	Password    string `json:"password"`
	OldPassword string `json:"old_password"`
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
}
