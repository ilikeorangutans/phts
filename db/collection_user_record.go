package db

// TODO  we're keeping this as a separate type because we want to add different ownership roles in the future
type CollectionUser struct {
	Timestamps
	CollectionID int64 `db:"collection_id"`
	UserID       int64 `db:"user_id"`
}
