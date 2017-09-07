package db

import "time"

type Clock func() time.Time

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
	Slug string `db:"slug"`
}

func (s *Sluggable) UpdateSlug(slug string) {
	s.Slug = slug
}
