package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

func AlbumListPhotosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)
	paginator := db.PaginatorFromRequest(r.URL.Query())

	db := model.DBFromRequest(r)
	backend := model.StorageFromRequest(r)

	colRepo := model.CollectionRepoFromRequest(r)
	applicableConfigs, err := colRepo.ApplicableRenditionConfigurations(collection)
	configs := RenditionConfigurationIDsFromQuery(applicableConfigs, r.URL.Query().Get("rendition-configuration-ids"))

	photoRepo := model.NewPhotoRepository(db, backend)

	album, _ := r.Context().Value("album").(model.Album)
	photos, paginator, err := photoRepo.ListAlbum(collection, album, paginator, configs)
	if err != nil {
		log.Fatal(err)
	}
	resp := ResponseWithPaginator{
		Data:      photos,
		Paginator: paginator,
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(resp)
	if err != nil {
		log.Fatal(err)
	}
}

func AlbumDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	album, _ := r.Context().Value("album").(model.Album)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(album)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)
	album, _ := r.Context().Value("album").(model.Album)

	db := model.DBFromRequest(r)
	albumRepo := model.NewAlbumRepository(db)

	err := albumRepo.Delete(collection, album)
	if err != nil {
		log.Fatal(err)
	}
}
func AddPhotosToAlbumHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(model.Collection)
	album, _ := r.Context().Value("album").(model.Album)
	db := model.DBFromRequest(r)
	albumRepo := model.NewAlbumRepository(db)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	submission := albumPhotoSubmission{}
	err := decoder.Decode(&submission)
	if err != nil {
		http.Error(w, "could not parse submitted json", http.StatusBadRequest)
		return
	}

	log.Printf("submission: %v", submission)
	log.Printf("album: %v", album)
	_, err = albumRepo.AddPhotos(collection, album, submission.PhotoIDs)
	if err != nil {
		log.Printf("error adding photos: %s", err.Error())
		http.Error(w, "could not add photos", http.StatusInternalServerError)
	}
}

type albumPhotoSubmission struct {
	AlbumID  int64   `json:"albumID"`
	PhotoIDs []int64 `json:"photoIDs"`
}

func RequireAlbum(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		collection, _ := r.Context().Value("collection").(model.Collection)
		albumID, err := strconv.ParseInt(chi.URLParam(r, "albumID"), 10, 64)
		if err != nil {
			panic(err)
		}

		db := model.DBFromRequest(r)
		albumRepo := model.NewAlbumRepository(db)
		album, err := albumRepo.FindByID(collection, albumID)
		if err != nil {
			log.Printf("error finding album: %s", err.Error())
			http.NotFound(w, r)
		}

		log.Printf("Found album %v", album)

		ctx := context.WithValue(r.Context(), "album", album)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
