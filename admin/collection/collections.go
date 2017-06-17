package collection

import (
	"context"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/web"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	foo, _ := r.Context().Value("foo").(string)
	log.Printf("Got foo: %s", foo)
	tmpl := web.GetTemplates("template/admin/base.tmpl", "template/admin/collection/index.tmpl")
	if tmpl == nil {
		return
	}

	colRepo := &model.DummyCollectionRepository{}
	col, _ := colRepo.FindByID(1)

	data := make(map[string]interface{})
	data["collection"] = col
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
	colRepo := &model.DummyCollectionRepository{}
	colRepo.Save(col)

	http.Redirect(w, r, "/admin/collections", http.StatusSeeOther)
}

func RequireCollection(wrap http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("requireCollectionFromSlug()")
		ctx := context.WithValue(r.Context(), "foo", "bar")
		r = r.WithContext(ctx)
		wrap(w, r)
	}
}
