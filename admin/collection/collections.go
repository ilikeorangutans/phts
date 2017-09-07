package collection

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
)

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	slug, _ := model.SlugFromString(r.PostFormValue("slug"))

	colRepo := model.CollectionRepoFromRequest(r)
	col := colRepo.Create(r.PostFormValue("name"), slug)
	var err error
	col, err = colRepo.Save(col)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/admin/collections", http.StatusSeeOther)
}

func UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	collection, _ := r.Context().Value("collection").(model.Collection)
	log.Printf("Uploading photo to %q", collection.Name)

	r.ParseMultipartForm(32 << 20)
	f, header, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	err = collection.AddPhoto(header.Filename, b)
	if err != nil {
		log.Panic(err)
	}
	repo := model.CollectionRepoFromRequest(r)
	repo.Save(collection)

	//http.Redirect(w, r, fmt.Sprintf("/admin/collections/%s", collection.Slug), http.StatusSeeOther)
}

func RequireCollection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		repo := model.CollectionRepoFromRequest(r)
		col, err := repo.FindBySlug(slug)
		if err != nil {
			log.Fatal(err)
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "collection", col)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func ServeRendition(w http.ResponseWriter, r *http.Request) {
	renditionID, err := strconv.Atoi(chi.URLParam(r, "rendition_id"))
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	dbx := model.DBFromRequest(r)
	repo := db.NewRenditionDB(dbx)

	rendition, err := repo.FindByID(int64(renditionID))
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	backend, _ := r.Context().Value("backend").(storage.Backend)
	data, err := backend.Get(rendition.ID)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}
