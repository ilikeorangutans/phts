package model

import "time"

type Timestamps struct {
	DateCreated time.Time
	DateUpdate  time.Time
}

type Record struct {
	ID     uint
	UserID uint
}

type Sluggable struct {
	Slug string
}

func (s *Sluggable) UpdateSlug(slug string) {
	s.Slug = slug
}

type User struct {
	Timestamps
	ID     uint
	Handle string
	Email  string
}

type Collection struct {
	Timestamps
	Record
	Sluggable

	newPhotos [][]byte
}

func (c *Collection) AddPhoto(data []byte) error {
	// TODO save the image data, start a job
	c.newPhotos = append(c.newPhotos, data)
	return nil
}

type Photo struct {
	Timestamps
	Record
	CollectionID uint
}

type Rendition struct {
	Timestamps
	Record
	PhotoID uint
}
