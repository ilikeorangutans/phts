package model

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/ilikeorangutans/phts/db"
)

type ShareSite struct {
	db.ShareSiteRecord
}

func (s ShareSite) Builder() ShareBuilder {
	return ShareBuilder{
		shareSite: s,
	}
}

type ShareBuilder struct {
	shareSite  ShareSite
	collection *db.Collection
	slug       string
	photos     []Photo
	errors     []error
	configs    RenditionConfigurations
}

func (b ShareBuilder) FromCollection(collection *db.Collection) ShareBuilder {
	b.collection = collection
	return b
}

func (b ShareBuilder) AddPhoto(photo Photo) ShareBuilder {
	b.photos = append(b.photos, photo)
	return b
}

func (b ShareBuilder) AllowRenditions(configs RenditionConfigurations) ShareBuilder {
	b.configs = configs
	return b
}

func (b ShareBuilder) WithSlug(input string) ShareBuilder {
	slug, err := SlugFromString(input)
	if err != nil {
		b.errors = append(b.errors, err)
	}
	b.slug = slug
	return b
}

func (b ShareBuilder) WithRandomSlug() ShareBuilder {
	input, err := randomHex(16)
	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	return b.WithSlug(input)
}

func (b ShareBuilder) Build() (Share, []error) {
	return Share{
		ShareSite:               b.shareSite,
		Photos:                  b.photos,
		Collection:              b.collection,
		RenditionConfigurations: b.configs,
		ShareRecord: db.ShareRecord{
			Slug: b.slug,
		},
	}, b.errors
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
