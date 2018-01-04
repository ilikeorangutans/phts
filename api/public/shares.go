package public

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/model"
)

func ViewShareHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := model.DBFromRequest(r)
	storage := model.StorageFromRequest(r)
	shareSite := r.Context().Value("shareSite").(model.ShareSite)
	shareRepo := model.NewShareRepository(db)
	photoRepo := model.NewPhotoRepository(db, storage)
	collectionRepo := model.NewCollectionRepository(db, storage)

	slug := chi.URLParam(r, "slug")

	share, err := shareRepo.FindByShareSiteAndSlug(shareSite, slug)
	if err != nil {
		log.Printf("No share found for slug %s and share site %s", slug, shareSite.Domain)
		http.NotFound(w, r)
		return
	}

	collection, err := collectionRepo.FindByID(share.CollectionID)
	if err != nil {
		log.Fatal(err)
	}
	photo, err := photoRepo.FindByID(collection, share.PhotoID)

	// TODO we're dumping the entire photo record, need something smaller here
	encoder := json.NewEncoder(w)
	resp := shareResponse{
		Share:  share,
		Photos: []model.Photo{photo},
	}
	err = encoder.Encode(resp)
	if err != nil {
		log.Fatal(err)
	}
}

func ServeShareRenditionHandler(w http.ResponseWriter, r *http.Request) {
	db := model.DBFromRequest(r)
	storage := model.StorageFromRequest(r)
	shareSite := r.Context().Value("shareSite").(model.ShareSite)
	shareRepo := model.NewShareRepository(db)
	collectionRepo := model.NewCollectionRepository(db, storage)
	renditionRepo := model.NewRenditionRepository(db)

	slug := chi.URLParam(r, "slug")

	share, err := shareRepo.FindByShareSiteAndSlug(shareSite, slug)
	if err != nil {
		log.Printf("No share found for slug %s and share site %s", slug, shareSite.Domain)
		http.NotFound(w, r)
		return
	}

	collection, err := collectionRepo.FindByID(share.CollectionID)
	if err != nil {
		log.Printf("could not find collection")
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "renditionID"), 10, 64)
	if err != nil {
		log.Printf("could not parse id")
		http.NotFound(w, r)
		return
	}
	rendition, err := renditionRepo.FindByID(collection, id)
	if err != nil {
		log.Printf("could not find rendition")
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Last-Modified", rendition.UpdatedAt.Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "max-age=3600")

	if modifiedSince := r.Header.Get("If-Modified-Since"); modifiedSince != "" {
		if timestamp, err := time.Parse(http.TimeFormat, modifiedSince); err == nil {
			if rendition.UpdatedAt.After(timestamp) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
	}

	data, err := storage.Get(id)
	if err != nil {
		log.Printf("binary not found for rendition: %v", err.Error())
		http.Error(w, "binary not found for rendition", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", rendition.Format)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))

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

type shareResponse struct {
	Share  model.Share   `json:"share"`
	Photos []model.Photo `json:"photos"`
}
