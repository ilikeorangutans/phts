package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	model2 "github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/web"

	"github.com/ilikeorangutans/phts/pkg/database"
	newmod "github.com/ilikeorangutans/phts/pkg/model"
)

type ResponseWithPaginator struct {
	Paginator database.Paginator `json:"paginator"`
	Data      interface{}        `json:"data"`
}

func PhotoStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dbx := web.DBFromRequest(r)
	user, err := web.UserFromRequest(r)
	if err != nil {
		http.Error(w, "no user", http.StatusInternalServerError)
		return
	}
	photoRepo := model2.NewPhotoRepo()

	paginator := database.PaginatorFromRequest(r.URL.Query())

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	photos, _, err := photoRepo.List(ctx, dbx, user, paginator)
	if err != nil {
		http.Error(w, "list failed", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(photos)
	if err != nil {
		http.Error(w, "encode failed", http.StatusInternalServerError)
		return
	}
}

func ListCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collectionRepo := model.CollectionRepoFromRequest(r)

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
	collection, _ := r.Context().Value("collection").(*db.Collection)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(collection)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(*db.Collection)

	colRepo := model.CollectionRepoFromRequest(r)
	err := colRepo.Delete(collection)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(collection)
	if err != nil {
		log.Fatal(err)
	}
}
func ServeRenditionHandler(w http.ResponseWriter, r *http.Request) {
	collection, _ := r.Context().Value("collection").(*db.Collection)

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

	data, err := backend.Get(id)
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

type createCollectionRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
	Slug string `json:"slug"`
}

// CreateCollectionFromRequest creates a createCollectionRequest from the request and sets some sensible defaults.
func CreateCollectionFromRequest(r *http.Request) (createCollectionRequest, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var createRequest createCollectionRequest
	err := decoder.Decode(&createRequest)
	if err != nil {
		return createRequest, errors.Wrap(err, "could not decode json")
	}

	createRequest.Name = strings.TrimSpace(createRequest.Name)
	createRequest.Slug = strings.TrimSpace(createRequest.Slug)

	if createRequest.Slug == "" {
		createRequest.Slug, err = model.SlugFromString(createRequest.Name)
		if err != nil {
			return createRequest, errors.Wrap(err, "could not create slug")
		}
	}
	return createRequest, nil
}

func CreateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	createRequest, err := CreateCollectionFromRequest(r)
	if err != nil {
		log.Printf("error creating request: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO here we'd do some validation

	owner, err := web.UserFromRequest(r)
	if err != nil {
		log.Println("no user in context")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	dbx := web.DBFromRequest(r)
	colRepo, err := newmod.NewCollectionRepo(dbx)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	collection, err := colRepo.NewCollection(ctx, createRequest.Name, createRequest.Slug, owner)
	if err != nil {
		log.Printf("error persisting: %s", err.Error())
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
	// TODO define error response format
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
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

	collection := web.CollectionFromRequest(r)

	dbx := web.DBFromRequest(r)
	storage := web.StorageBackendFromRequest(r)
	collectionRepo, _ := model2.NewCollectionRepo(dbx)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	photoUpload, err := model2.FromReader(file, fileHeader.Filename)
	if err != nil {
		log.Printf("error creating upload from request file %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel = context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	queue := web.GetRenditionUpdateRequestQueueFromRequest(r)

	collection, photos, err := collectionRepo.AddPhotos(ctx, dbx, storage, collection, queue, photoUpload)
	if err != nil {
		log.Printf("could not add photo: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(photos[0])
	if err != nil {
		log.Fatal(err)
	}
}

func RenditionConfigurationIDsFromQuery(applicableConfigs model.RenditionConfigurations, configIDs string) (result []model.RenditionConfiguration) {
	if len(configIDs) > 0 {
		split := strings.Split(configIDs, ",")

		for _, s := range split {
			id, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Printf("failed to parse rendition configuration id %s", err.Error())
				continue
			}
			config, err := applicableConfigs.ByID(id)
			if err != nil {
				log.Printf("no config: %s", err.Error())
				continue
			}

			result = append(result, config)
		}
	}

	if len(result) == 0 {
		result = applicableConfigs
	}

	return result
}

func ListRecentPhotosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(*db.Collection)

	var err error
	colRepo := model.CollectionRepoFromRequest(r)
	applicableConfigs, err := colRepo.ApplicableRenditionConfigurations(collection)
	configs := RenditionConfigurationIDsFromQuery(applicableConfigs, r.URL.Query().Get("rendition-configuration-ids"))

	photoRepo := model.PhotoRepoFromRequest(r)
	photos, paginator, err := photoRepo.List(collection, database.PaginatorFromRequest(r.URL.Query()), configs)
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

func CreatePhotoShareHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(*db.Collection)

	db := model.DBFromRequest(r)
	storage := model.StorageFromRequest(r)
	shareRepo := model.NewShareRepository(db)
	shareSiteRepo := model.NewShareSiteRepository(db)
	photoRepo := model.NewPhotoRepository(db, storage)
	collectionRepo := model.CollectionRepoFromRequest(r)
	renditionConfigs, err := collectionRepo.ApplicableRenditionConfigurations(collection)
	if err != nil {
		log.Printf("error loading rendition configs: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shareRequest, err := ShareRequestFromRequest(r)
	if err != nil {
		log.Printf("error parsing JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shareSite, err := shareSiteRepo.FindByID(shareRequest.ShareSiteID)
	if err != nil {
		log.Printf("cannot find share site: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	photo, err := photoRepo.FindByID(collection, shareRequest.PhotoID)
	if err != nil {
		log.Printf("cannot find photo: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	builder := shareSite.Builder().
		FromCollection(collection).
		AddPhoto(photo).
		AllowRenditions(shareRequest.FilterRenditionConfigurations(renditionConfigs))

	if shareRequest.GenerateRandomSlug() {
		builder = builder.WithRandomSlug()
	} else {
		builder = builder.WithSlug(shareRequest.Slug)
	}

	share, errors := builder.Build()
	if len(errors) > 0 {
		log.Printf("errors from builder: %v", errors)
		http.Error(w, "error from builder", http.StatusBadRequest)
		return
	}
	share, err = shareRepo.Publish(share)
	if err != nil {
		log.Printf("error saving: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(shareRequest)
}

func ShowPhotoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(*db.Collection)

	colRepo := model.CollectionRepoFromRequest(r)
	applicableConfigs, err := colRepo.ApplicableRenditionConfigurations(collection)
	configs := RenditionConfigurationIDsFromQuery(applicableConfigs, r.URL.Query().Get("rendition-configuration-ids"))

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

	renditionRepo := model.NewRenditionRepository(db)
	renditions, err := renditionRepo.FindByPhotoAndRenditionConfigurations(collection, photo, configs)
	if err != nil {
		log.Printf("could not load renditions", err.Error())
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	photo.Renditions = renditions
	w.Header().Set("Last-Modified", photo.UpdatedAt.Format(http.TimeFormat))

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(photo); err != nil {
		log.Fatal(err)
	}
}

func DeletePhotoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(*db.Collection)

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

	err = photoRepo.Delete(collection, photo)

	if err != nil {
		log.Printf("Error deleting photo %s", err.Error())
		http.Error(w, `{"message":"error deleting photo"}`, http.StatusInternalServerError)
	}
}
func ShowPhotoSharesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(*db.Collection)

	dbx := model.DBFromRequest(r)
	backend := model.StorageFromRequest(r)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	shareRepo := model.NewShareRepository(dbx)
	photoRepo := model.NewPhotoRepository(dbx, backend)
	photo, err := photoRepo.FindByID(collection, id)
	if err != nil {
		log.Printf("photo not found: %v", err.Error())
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	paginator := database.PaginatorFromRequest(r.URL.Query())

	shares, err := shareRepo.FindByPhoto(photo, paginator)
	if err != nil {
		log.Printf("shares not found: %v", err.Error())
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(shares); err != nil {
		log.Fatal(err)
	}
}
func ListRenditionConfigurationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collectionRepo := model.CollectionRepoFromRequest(r)

	collection, _ := r.Context().Value("collection").(*db.Collection)
	configs, err := collectionRepo.ApplicableRenditionConfigurations(collection)
	if err != nil {
		log.Printf("error retrieving configurations: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO we don't really have paginator support yet.
	paginator := database.PaginatorFromRequest(r.URL.Query())
	withPaginator := ResponseWithPaginator{
		Paginator: paginator,
		Data:      configs,
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(withPaginator)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateRenditionConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbx := model.DBFromRequest(r)
	repo := model.NewRenditionConfigurationRepository(dbx)

	collection, _ := r.Context().Value("collection").(*db.Collection)

	config := model.RenditionConfiguration{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&config)
	if err != nil {
		log.Printf("error parsing JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	config, err = repo.Save(collection, config)
	if err != nil {
		log.Printf("error saving: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(config)
}

func ListPhotosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection, _ := r.Context().Value("collection").(*db.Collection)

	paginator := database.PaginatorFromRequest(r.URL.Query())

	db := model.DBFromRequest(r)
	backend := model.StorageFromRequest(r)

	colRepo := model.CollectionRepoFromRequest(r)
	applicableConfigs, err := colRepo.ApplicableRenditionConfigurations(collection)
	configs := RenditionConfigurationIDsFromQuery(applicableConfigs, r.URL.Query().Get("rendition-configuration-ids"))

	photoRepo := model.NewPhotoRepository(db, backend)
	photos, paginator, err := photoRepo.List(collection, paginator, configs)
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

func ListShareSitesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO records need to be scoped here

	db := model.DBFromRequest(r)

	repo := model.NewShareSiteRepository(db)

	sites, err := repo.List()
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(sites)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateShareSitesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := model.DBFromRequest(r)

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var shareSite model.ShareSite
	err := decoder.Decode(&shareSite)
	if err != nil {
		log.Printf("error parsing JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shareSiteRepo := model.NewShareSiteRepository(db)
	shareSite, err = shareSiteRepo.Save(shareSite)
	if err != nil {
		log.Printf("error parsing JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(shareSite)
	if err != nil {
		log.Fatal(err)
	}
}
