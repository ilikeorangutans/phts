package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
)

func CreateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var album model.Album
	err := decoder.Decode(&album)
	if err != nil {
		log.Printf("error parsing JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	album.CollectionID = collection.ID

	db := model.DBFromRequest(r)
	repo := model.NewAlbumRepository(db)
	album, err = repo.Save(album)
	if err != nil {
		log.Printf("error parsing JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(album)
	if err != nil {
		log.Fatal(err)
	}
}

func ListAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)

	paginator := db.PaginatorFromRequest(r.URL.Query())

	db := model.DBFromRequest(r)
	albumRepo := model.NewAlbumRepository(db)
	albums, paginator, err := albumRepo.List(collection, paginator)
	if err != nil {
		log.Fatal(err)
	}

	resp := ResponseWithPaginator{
		Data:      albums,
		Paginator: paginator,
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(resp)
	if err != nil {
		log.Fatal(err)
	}
}
