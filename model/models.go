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

type Collection struct {
	Timestamps
	Record
	Sluggable

	Name string `db:"name"`
}

func (c *Collection) AddPhoto(data []byte) (Collection, error) {
	p := Photo{
		Timestamps:   JustCreated(),
		CollectionID: c.ID,
		Renditions: []Rendition{
			Rendition{
				Timestamps: JustCreated(),
				Data:       data,
			},
		},
	}

	c.Photos = append(c.Photos, p)

	return *c, nil
}

type Photo struct {
	Timestamps
	Record
	Renditions   []Rendition
	CollectionID int64  `db:"collection_id"`
	Description  string `db:"description"`
}

type Rendition struct {
	Timestamps
	Record
	Original bool
	PhotoID  int64
	Data     []byte
}
