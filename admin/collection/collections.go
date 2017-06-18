package collection

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/web"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := web.GetTemplates("template/admin/base.tmpl", "template/admin/collection/index.tmpl")
	if tmpl == nil {
		return
	}

	colRepo := model.CollectionRepoFromRequest(r)
	cols, err := colRepo.Recent(10)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]interface{})
	data["collections"] = cols
	tmpl.Execute(w, data)
}

func NewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := web.GetTemplates("template/admin/base.tmpl", "template/admin/collection/new.tmpl")

	data := make(map[string]interface{})
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	slug, _ := model.SlugFromString(r.PostFormValue("slug"))

	col := model.Collection{
		Timestamps: model.JustCreated(),
		Name:       r.PostFormValue("name"),
		Sluggable:  model.Sluggable{Slug: slug},
	}
	colRepo := model.CollectionRepoFromRequest(r)
	var err error
	col, err = colRepo.Save(col)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/admin/collections", http.StatusSeeOther)
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	collection, _ := r.Context().Value("collection").(model.Collection)

	tmpl := web.GetTemplates("template/admin/base.tmpl", "template/admin/collection/show.tmpl")
	data := make(map[string]interface{})
	data["collection"] = collection
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	collection, _ := r.Context().Value("collection").(model.Collection)
	log.Printf("Uploading photo to %s", collection)

	http.Redirect(w, r, fmt.Sprintf("/admin/collections/%s", collection.Slug), http.StatusSeeOther)
}

func RequireCollection(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug, ok := vars["slug"]
		if !ok {
			http.NotFound(w, r)
			return
		}

		repo := model.CollectionRepoFromRequest(r)
		col, err := repo.FindBySlug(slug)
		if err != nil {
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "collection", col)
		r = r.WithContext(ctx)
		wrap(w, r)
	}
}
