package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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
