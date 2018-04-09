package public

import (
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
)

type ShareRepository interface {
	FindShareBySlug(model.ShareSite, string) (viewShareResponse, error)
}

func NewShareRepository(dbx db.DB, collectionRepo model.CollectionRepository, storage storage.Backend) ShareRepository {
	return &shareRepo{
		shareRepo:      model.NewShareRepository(dbx),
		collectionRepo: collectionRepo,
		photoRepo:      model.NewPhotoRepository(dbx, storage),
	}
}

type shareRepo struct {
	shareRepo      model.ShareRepository
	collectionRepo CollectionRepository
	photoRepo      model.PhotoRepository
}

func (r *shareRepo) FindShareBySlug(shareSite model.ShareSite, slug string) (response viewShareResponse, err error) {
	share, err := r.shareRepo.FindByShareSiteAndSlug(shareSite, slug)
	if err != nil {
		return response, err
	}

	collection, err := r.collectionRepo.FindByID(share.CollectionID)
	if err != nil {
		log.Fatal(err)
	}
	renditionConfigs := []model.RenditionConfiguration{}
	sharedConfigs := []sharedRenditionConfiguration{}
	for _, c := range renditionConfigs {
		sharedConfigs = append(sharedConfigs, newSharedRenditionConfiguration(c))
	}

	photo, err := r.photoRepo.FindByID(collection, share.PhotoID)
	photos := []sharedPhoto{
		newSharedPhoto(photo),
	}

	response = viewShareResponse{
		Share:                   shareResponse{Slug: share.Slug},
		Photos:                  photos,
		RenditionConfigurations: sharedConfigs,
	}

	return response, nil
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
