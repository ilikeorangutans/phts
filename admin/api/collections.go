package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
)

type ResponseWithPaginator struct {
	Paginator db.Paginator `json:"paginator"`
	Data      interface{}  `json:"data"`
}

func ListCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := model.DBFromRequest(r)
	backend := model.StorageFromRequest(r)

	collectionRepo := model.NewCollectionRepository(db, backend)

	collections, err := collectionRepo.Recent(100)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(collections)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(collection)
	if err != nil {
		log.Fatal(err)
	}
}

func ServeRenditionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)

	db := model.DBFromRequest(r)
	backend := model.StorageFromRequest(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	renditionRepo := model.NewRenditionRepository(db)
	rendition, err := renditionRepo.FindByID(collection, id)
	if err != nil {
		log.Printf("rendition not found: %v", err.Error())
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	data, err := backend.Get(id)
	if err != nil {
		log.Printf("binary not found for rendition: %v", err.Error())
		http.Error(w, "binary not found for rendition", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", rendition.Format)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Last-Modified", rendition.UpdatedAt.Format(http.TimeFormat))

	if r.Method == "HEAD" {
		return
	}

	written, err := io.Copy(w, bytes.NewReader(data))
	if err != nil {
		log.Fatalf("error while writing binary to respose: %s", err.Error())
	}
	if written < int64(len(data)) {
		log.Printf("wrote %d/%d bytes", written, len(data))
	}
}

func CreateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var collection model.Collection
	err := decoder.Decode(&collection)
	if err != nil {
		log.Printf("error parsing JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO here we'd do some validation

	colRepo := model.CollectionRepoFromRequest(r)
	collection, err = colRepo.Save(collection)
	if err != nil {
		log.Printf("error parsing JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(collection)
	if err != nil {
		log.Fatal(err)
	}
}

func UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseMultipartForm(32 << 23)
	if err != nil {
		log.Printf("error parsing form: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Printf("error parsing form: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("uploaded file: %s", fileHeader.Filename)

	data, err := ioutil.ReadAll(file)

	log.Printf("file size: %d", len(data))

	collection, _ := r.Context().Value("collection").(model.Collection)
	colRepo := model.CollectionRepoFromRequest(r)

	photo, err := colRepo.AddPhoto(collection, fileHeader.Filename, data)
	if err != nil {
		log.Printf("error parsing form: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(photo)
	if err != nil {
		log.Fatal(err)
	}
}

func ListRecentPhotosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)

	colRepo := model.CollectionRepoFromRequest(r)
	configs, err := colRepo.ApplicableRenditionConfigurations(collection)
	if err != nil {
		log.Fatal(err)
	}
	photoRepo := model.PhotoRepoFromRequest(r)
	photos, paginator, err := photoRepo.List(collection, db.PaginatorFromRequest(r.URL.Query()), configs)
	if err != nil {
		log.Fatal(err)
	}

	withPaginator := ResponseWithPaginator{
		Paginator: paginator,
		Data:      photos,
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(withPaginator)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowPhotoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)

	db := model.DBFromRequest(r)
	backend := model.StorageFromRequest(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	photoRepo := model.NewPhotoRepository(db, backend)
	photo, err := photoRepo.FindByID(collection, id)
	if err != nil {
		log.Printf("photo not found: %v", err.Error())
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Last-Modified", photo.UpdatedAt.Format(http.TimeFormat))

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(photo); err != nil {
		log.Fatal(err)
	}
}
