package model

import "github.com/ilikeorangutans/phts/db"

type AlbumRepository interface {
	Save(Album) (Album, error)
	List(Collection, db.Paginator) ([]Album, db.Paginator, error)
}

func NewAlbumRepository(dbx db.DB) AlbumRepository {
	return &albumRepoImpl{
		albumDB: db.NewAlbumDB(dbx),
	}
}

type albumRepoImpl struct {
	albumDB db.AlbumDB
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
