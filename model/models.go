package model

import "time"

type Timestamps struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func JustCreated() Timestamps {
	return Timestamps{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type Record struct {
	ID     int64 `db:"id"`
	UserID int64
}

type Sluggable struct {
	Slug string `db:"slug"`
}

func (s *Sluggable) UpdateSlug(slug string) {
	s.Slug = slug
}

type Collection struct {
	Timestamps
	Record
	Sluggable

	Name string `db:"name"`

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
	CollectionID int64  `db:"collection_id"`
	Description  string `db:"description"`
}

type Rendition struct {
	Timestamps
	Record
	PhotoID int64
}
