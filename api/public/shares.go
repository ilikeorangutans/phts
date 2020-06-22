package public

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/ilikeorangutans/phts/model"
	newmodel "github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/web"
)

func ViewShareHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	slug := chi.URLParam(r, "slug")
	dbx := web.DBFromRequest(r)
	shareSite := r.Context().Value(web.ShareSiteKey).(newmodel.ShareSite)
	share, err := newmodel.FindSharedPhotoBySlug(ctx, dbx, shareSite, slug)
	if err != nil {
		log.Printf("could not get share: %v", err)
		http.NotFound(w, r)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(newViewShareResponse(share))
	if err != nil {
		log.Fatal(err)
	}
}

func newViewShareResponse(share newmodel.ShareWithPhotos) viewShareResponse {
	var photos []sharedPhoto
	for _, photo := range share.Photos {
		var renditions []sharedRendition
		for _, rendition := range photo.Renditions {
			renditions = append(renditions, sharedRendition{
				ID:                       rendition.ID,
				Width:                    rendition.Width,
				Height:                   rendition.Height,
				Original:                 rendition.Original,
				RenditionConfigurationID: rendition.RenditionConfigurationID,
			})
		}

		photos = append(photos, sharedPhoto{
			ID:         photo.Photo.ID,
			Renditions: renditions,
		})
	}

	var renditions []sharedRenditionConfiguration
	for _, config := range share.RenditionConfigurations {
		renditions = append(renditions, sharedRenditionConfiguration{
			ID:       config.ID,
			Width:    config.Width,
			Height:   config.Height,
			Original: config.Original,
		})
	}

	return viewShareResponse{
		Share: shareResponse{
			Slug:      share.Share.Slug,
			CreatedAt: share.Share.CreatedAt,
		},
		Photos:                  photos,
		RenditionConfigurations: renditions,
	}
}

type viewShareResponse struct {
	Share                   shareResponse                  `json:"share"`
	Photos                  []sharedPhoto                  `json:"photos"`
	RenditionConfigurations []sharedRenditionConfiguration `json:"rendition_configurations"`
}

type sharedRenditionConfiguration struct {
	ID       int64 `json:"id"`
	Width    int   `json:"width"`
	Height   int   `json:"height"`
	Original bool  `json:"original"`
}

type shareResponse struct {
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type sharedPhoto struct {
	ID         int64             `json:"id"`
	Renditions []sharedRendition `json:"renditions"`
}

type sharedRendition struct {
	ID                       int64 `json:"id"`
	Width                    uint  `json:"width"`
	Height                   uint  `json:"height"`
	Original                 bool  `json:"original"`
	RenditionConfigurationID int64 `json:"rendition_configuration_id"`
}

func ServeShareRenditionHandler(w http.ResponseWriter, r *http.Request) {
	db := model.DBFromRequest(r)
	storage := model.StorageFromRequest(r)
	shareSite := r.Context().Value("shareSite").(model.ShareSite)
	shareRepo := model.NewShareRepository(db)
	renditionRepo := model.NewRenditionRepository(db)

	slug := chi.URLParam(r, "slug")

	share, err := shareRepo.FindByShareSiteAndSlug(shareSite, slug)
	if err != nil {
		log.Printf("No share found for slug %s and share site %s", slug, shareSite.Domain)
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "renditionID"), 10, 64)
	if err != nil {
		log.Printf("could not parse id")
		http.NotFound(w, r)
		return
	}

	rendition, err := renditionRepo.FindByShareAndID(share, id)
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
