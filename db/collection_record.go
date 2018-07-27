package db

import (
	"log"
	"time"
)

// Collection is a single database level record of a collection.
type Collection struct {
	Record
	Timestamps
	Sluggable
	Name       string `db:"name" json:"name"`
	PhotoCount int    `db:"photo_count" json:"photoCount"`
}

type CollectionDB interface {
	FindByID(id int64) (Collection, error)
	FindBySlug(slug string) (Collection, error)
	Save(collection Collection) (Collection, error)
	List(userID int64, count int, afterID int64, orderBy string) ([]Collection, error)
	Delete(int64) error
	Assign(userID int64, collectionID int64) error
	CanAccess(userID int64, collectionID int64) bool
}

func NewCollectionDB(db DB) CollectionDB {
	return &collectionSQLDB{
		clock: time.Now,
		db:    db,
	}
}

type collectionSQLDB struct {
	db    DB
	clock Clock
}

func (c *collectionSQLDB) CanAccess(userID int64, collectionID int64) bool {
	// TODO implement me
	return true
}

func (c *collectionSQLDB) List(userID int64, count int, afterID int64, orderBy string) ([]Collection, error) {
	result := []Collection{}
	sql := "SELECT c.* FROM collections AS c, users_collections AS uc WHERE uc.user_id = $1 AND uc.collection_id = c.id AND c.id > $2 order by $3 limit $4"
	rows, err := c.db.Queryx(
		sql,
		userID,
		afterID,
		orderBy,
		count,
	)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		record := Collection{}
		err := rows.StructScan(&record)
		if err != nil {
			return result, err
		}
		result = append(result, record)
	}

	return result, nil
}

func (c *collectionSQLDB) Assign(userID int64, collectionID int64) error {
	log.Printf("Assigning collection %d to user %d", collectionID, userID)
	_, err := c.db.Exec("INSERT INTO users_collections (user_id, collection_id) VALUES ($1, $2)", userID, collectionID)
	log.Println(err)
	return err
}

func (c *collectionSQLDB) FindByID(id int64) (Collection, error) {
	var record Collection
	err := c.db.QueryRowx("SELECT * FROM collections WHERE id = $1 LIMIT 1", id).StructScan(&record)
	return record, err
}

func (c *collectionSQLDB) Delete(id int64) error {
	sql := "DELETE FROM collections WHERE id=$1"
	return checkResult(c.db.Exec(sql, id))
}

func (c *collectionSQLDB) FindBySlug(slug string) (Collection, error) {
	var record Collection
	err := c.db.QueryRowx("SELECT * FROM collections WHERE slug = $1 LIMIT 1", slug).StructScan(&record)
	return record, err
}

func (c *collectionSQLDB) Save(record Collection) (Collection, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(c.clock)
		sql := "UPDATE collections SET name = $1, slug = $2, updated_at = $3, photo_count = (SELECT count(*) FROM photos WHERE collection_id = $4) WHERE id = $4"
		record.UpdatedAt = c.clock()
		err = checkResult(c.db.Exec(
			sql,
			record.Name,
			record.Slug,
			record.UpdatedAt,
			record.ID,
		))
	} else {
		record.Timestamps = JustCreated(c.clock)
		sql := "INSERT INTO collections (name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
		err = c.db.QueryRow(sql, record.Name, record.Slug, record.CreatedAt, record.UpdatedAt).Scan(&record.ID)
	}

	return record, err
}
