package photo

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/web"
)

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	collection := r.Context().Value("collection").(model.Collection)
	photo := r.Context().Value("photo").(model.Photo)

	tmpl := web.GetTemplates("template/admin/base.tmpl", "template/admin/collection/photo/show.tmpl")
	data := make(map[string]interface{})
	data["photo"] = photo
	data["collection"] = collection

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func RequirePhoto(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		photoIdString, ok := vars["photo_id"]
		if !ok {
			log.Printf("Photo with id %s not found", photoIdString)
			http.NotFound(w, r)
			return
		}

		photoID, err := strconv.ParseInt(photoIdString, 10, 64)
		if err != nil {
			log.Printf("Invalid photo ID %q", photoIdString)
			http.NotFound(w, r)
			return
		}

		collection := r.Context().Value("collection").(model.Collection)

		repo := model.PhotoRepoFromRequest(r)
		photo, err := repo.FindByID(collection.ID, photoID)
		if err != nil {
			log.Panic(err)
		}

		ctx := context.WithValue(r.Context(), "photo", photo)
		r = r.WithContext(ctx)
		wrap(w, r)
	}
}
