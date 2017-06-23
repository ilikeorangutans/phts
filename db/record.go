package db

import "time"

type Record struct {
	ID int64 `db:"id"`
}

func (r Record) IsPersisted() bool {
	return r.ID != 0
}

type Timestamps struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (t *Timestamps) JustUpdated() {
	t.UpdatedAt = time.Now()
}

func JustCreated() Timestamps {
	return Timestamps{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type Sluggable struct {
	Slug string `db:"slug"`
}

func (s *Sluggable) UpdateSlug(slug string) {
	s.Slug = slug
}
