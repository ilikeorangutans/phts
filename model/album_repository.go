package model

import (
	"log"

	"github.com/ilikeorangutans/phts/db"
)

type AlbumRepository interface {
	Save(Album) (Album, error)
	List(Collection, db.Paginator) ([]Album, db.Paginator, error)
	FindByID(Collection, int64) (Album, error)
	AddPhotos(Collection, Album, []int64) (Album, error)
	Delete(Collection, Album) error
}

func NewAlbumRepository(dbx db.DB) AlbumRepository {
	return &albumRepoImpl{
		albumDB: db.NewAlbumDB(dbx),
	}
}

type albumRepoImpl struct {
	albumDB db.AlbumDB
}

func (r *albumRepoImpl) FindByID(collection Collection, id int64) (Album, error) {
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

func (r *albumRepoImpl) List(collection Collection, paginator db.Paginator) ([]Album, db.Paginator, error) {
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

func (r *albumRepoImpl) AddPhotos(collection Collection, album Album, photoIDs []int64) (Album, error) {
	log.Printf("Adding photos %v", photoIDs)
	err := r.albumDB.AddPhotos(collection.ID, album.ID, photoIDs)
	return Album{}, err
}

func (r *albumRepoImpl) Delete(collection Collection, album Album) error {
	return r.albumDB.Delete(collection.ID, album.ID)
}
