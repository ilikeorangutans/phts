package model

type Record struct {
	ID int64 `db:"id"`
}

func (r Record) IsPersisted() bool {
	return r.ID != 0
}
