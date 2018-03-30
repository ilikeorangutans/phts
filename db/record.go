package db

import "time"

type Clock func() time.Time

type X struct {
	Record
	Timestamps
}

type Record struct {
	ID int64 `db:"id" json:"id"`
}

func (r Record) IsPersisted() bool {
	return r.ID != 0
}

type Timestamps struct {
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func (t *Timestamps) JustUpdated(clock Clock) {
	t.UpdatedAt = clock()
}

func JustCreated(clock Clock) Timestamps {
	return Timestamps{
		CreatedAt: clock(),
		UpdatedAt: clock(),
	}
}

type Sluggable struct {
	Slug string `db:"slug" json:"slug"`
}

func (s *Sluggable) UpdateSlug(slug string) {
	s.Slug = slug
}
