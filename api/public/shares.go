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
	renditionConfigs, err := collectionRepo.ApplicableRenditionConfigurations(collection)
	if err != nil {
		log.Println("no rendition configurations found")
		http.NotFound(w, r)
		return
	}

	sharedConfigs := []sharedRenditionConfiguration{}
	for _, c := range renditionConfigs {
		sharedConfigs = append(sharedConfigs, newSharedRenditionConfiguration(c))
	}

	photo, err := photoRepo.FindByID(collection, share.PhotoID)
	photos := []sharedPhoto{
		newSharedPhoto(photo),
	}

	encoder := json.NewEncoder(w)
	resp := viewShareResponse{
		Share:                   shareResponse{Slug: share.Slug},
		Photos:                  photos,
		RenditionConfigurations: sharedConfigs,
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

type viewShareResponse struct {
	Share                   shareResponse                  `json:"share"`
	Photos                  []sharedPhoto                  `json:"photos"`
	RenditionConfigurations []sharedRenditionConfiguration `json:"rendition_configurations"`
}

type shareResponse struct {
	Slug string `json:"slug"`
}

func newSharedPhoto(photo model.Photo) sharedPhoto {
	renditions := []sharedRendition{}
	for _, r := range photo.Renditions {
		renditions = append(renditions, sharedRendition{
			ID:     r.ID,
			Width:  r.Width,
			Height: r.Height,
			RenditionConfigurationID: r.RenditionConfigurationID,
		})
	}

	return sharedPhoto{
		Renditions: renditions,
	}
}

type sharedPhoto struct {
	Renditions []sharedRendition `json:"renditions"`
}

type sharedRendition struct {
	ID                       int64 `json:"id"`
	Width                    uint  `json:"width"`
	Height                   uint  `json:"height"`
	RenditionConfigurationID int64 `json:"rendition_configuration_id"`
}

func newSharedRenditionConfiguration(config model.RenditionConfiguration) sharedRenditionConfiguration {
	return sharedRenditionConfiguration{
		ID:       config.ID,
		Width:    config.Width,
		Height:   config.Height,
		Original: config.Original,
	}
}

type sharedRenditionConfiguration struct {
	ID       int64 `json:"id"`
	Width    int   `json:"width"`
	Height   int   `json:"height"`
	Original bool  `json:"original"`
}
