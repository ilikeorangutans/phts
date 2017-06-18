package model

import (
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type CollectionRepository interface {
	FindByID(id uint) (Collection, error)
	FindBySlug(slug string) (Collection, error)
	Save(Collection) (Collection, error)
	Recent(int) ([]Collection, error)
}

type CollectionSQLRepository struct {
	db *sqlx.DB
}

func (r *CollectionSQLRepository) FindBySlug(slug string) (Collection, error) {
	var result Collection
	err := r.db.QueryRowx("SELECT * FROM collections WHERE slug=$1", slug).StructScan(&result)
	return result, err
}

func (r *CollectionSQLRepository) FindByID(id uint) (Collection, error) {
	var result Collection
	err := r.db.QueryRow("SELECT * FROM collections WHERE id=?", id).Scan(&result)
	return result, err
}

func (r *CollectionSQLRepository) Recent(n int) ([]Collection, error) {

	rows, err := r.db.Queryx("SELECT * FROM collections ORDER BY updated_at DESC LIMIT $1", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Collection{}
	for rows.Next() {
		col := Collection{}
		err = rows.StructScan(&col)
		if err != nil {
			return nil, err
		}
		result = append(result, col)
	}

	return result, nil
}

func (r *CollectionSQLRepository) Save(col Collection) (Collection, error) {
	if col.ID != 0 {
		log.Printf("Collection already saved, doing nothing")
		return col, nil
	}

	now := time.Now()
	err := r.db.QueryRow("INSERT INTO collections (name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", col.Name, col.Slug, now, now).Scan(&col.ID)
	if err != nil {
		return col, err
	}

	return col, nil
}

func DBFromRequest(r *http.Request) *sqlx.DB {
	db, ok := r.Context().Value("database").(*sqlx.DB)
	if !ok {
		log.Fatal("Could not get database from request, wrong type")
	}

	return db
}

func CollectionRepoFromRequest(r *http.Request) CollectionRepository {
	return &CollectionSQLRepository{
		db: DBFromRequest(r),
	}
}
