package api

import (
	"encoding/json"
	"net/http"

	"github.com/ilikeorangutans/phts/model"
)

func ShareRequestFromRequest(r *http.Request) (result ShareRequest, err error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&result)
	return result, err
}

type ShareRequest struct {
	PhotoID           int64   `json:"photoID"`
	ShareSiteID       int64   `json:"shareSiteID"`
	SlugStrategy      string  `json:"slugStrategy"`
	Slug              string  `json:"slug"`
	AllowedRenditions []int64 `json:"allowedRenditions"`
}

func (s ShareRequest) GenerateRandomSlug() bool {
	return s.SlugStrategy == "random"
}

func (s ShareRequest) FilterRenditionConfigurations(input model.RenditionConfigurations) model.RenditionConfigurations {
	result := model.RenditionConfigurations{}
	if len(s.AllowedRenditions) == 0 {
		// TODO not sure if this is a good default
		return input
	}

	for _, config := range input {
		for _, id := range s.AllowedRenditions {
			if config.ID == id {
				result = append(result, config)
			}
		}
	}

	return result
}
