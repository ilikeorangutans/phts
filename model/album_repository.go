package model

import (
	"log"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
)

type AlbumRepository interface {
	Save(Album) (Album, error)
	List(db.Collection, database.Paginator) ([]Album, database.Paginator, error)
	FindByID(db.Collection, int64) (Album, error)
	AddPhotos(db.Collection, Album, []int64) (Album, error)
	Delete(db.Collection, Album) error
}

func NewAlbumRepository(dbx db.DB) AlbumRepository {
	return &albumRepoImpl{
		albumDB: db.NewAlbumDB(dbx),
	}
}

type albumRepoImpl struct {
	albumDB db.AlbumDB
}

func (r *albumRepoImpl) FindByID(collection db.Collection, id int64) (Album, error) {
	record, err := r.albumDB.FindByID(collection.ID, id)
	if err != nil {
		return Album{}, err
	}

	album := Album{record}

	return album, nil
}

func (r *albumRepoImpl) Save(album Album) (Album, error) {
	if len(album.Slug) == 0 {
		slug, err := SlugFromString(album.Name)
		if err != nil {
			return album, err
		}

		album.Slug = slug
	}
	record, err := r.albumDB.Save(album.AlbumRecord)
	album.AlbumRecord = record
	return album, err
}

func (r *albumRepoImpl) List(collection db.Collection, paginator database.Paginator) ([]Album, database.Paginator, error) {
	records, err := r.albumDB.List(collection.ID, paginator)
	if err != nil {
		return nil, paginator, err
	}

	result := []Album{}
	for _, record := range records {
		result = append(result, Album{record})
	}

	return result, paginator, nil
}

func (r *albumRepoImpl) AddPhotos(collection db.Collection, album Album, photoIDs []int64) (Album, error) {
	log.Printf("Adding photos %v", photoIDs)
	err := r.albumDB.AddPhotos(collection.ID, album.ID, photoIDs)
	return Album{}, err
}

func (r *albumRepoImpl) Delete(collection db.Collection, album Album) error {
	return r.albumDB.Delete(collection.ID, album.ID)
}
