package collection

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/ilikeorangutans/phts/web"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := web.GetTemplates("template/admin/base.tmpl", "template/admin/collection/index.tmpl")
	if tmpl == nil {
		return
	}

	colRepo := model.CollectionRepoFromRequest(r)
	cols, err := colRepo.Recent(12)
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

func SaveHandlerJSON(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var collection model.Collection
	decoder.Decode(&collection)

	colRepo := model.CollectionRepoFromRequest(r)

	res, err := colRepo.Save(collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		log.Println("Successfully created ", res.ID, res.Name, res.CreatedAt)
	}
}

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

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	collection, _ := r.Context().Value("collection").(model.Collection)
	repo := model.CollectionRepoFromRequest(r)
	photos, err := repo.RecentPhotos(collection, 12)
	if err != nil {
		// TODO return 404?
		log.Println(err)
	}

	tmpl := web.GetTemplates("template/admin/base.tmpl", "template/admin/collection/show.tmpl", "template/admin/collection/photo/card.tmpl")
	data := make(map[string]interface{})
	data["collection"] = collection
	data["recentPhotos"] = photos
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func ListPhotosHandler(w http.ResponseWriter, r *http.Request) {
	collection, _ := r.Context().Value("collection").(model.Collection)
	repo := model.PhotoRepoFromRequest(r)

	paginator := db.PaginatorFromRequest(r.URL.Query())
	log.Printf("Got paginator %v", paginator)
	paginator.Count = 12
	photos, paginator, err := repo.List(collection, paginator)
	if err != nil {
		log.Panic(err)
	}

	data := make(map[string]interface{})
	data["collection"] = collection
	data["photos"] = photos
	data["paginator"] = paginator
	// TODO need a way to determine if there's more photos
	// Also need a way to find photos "before" ID, aka paginate backwards

	tmpl := web.GetTemplates("template/admin/base.tmpl", "template/admin/collection/show.tmpl", "template/admin/collection/photo/list.tmpl", "template/admin/collection/photo/card.tmpl")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
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
