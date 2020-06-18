package model

import (
	"context"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
)

func NewPhotoRepo() *PhotoRepo {
	return &PhotoRepo{
		clock: time.Now,
		stmt:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type PhotoRepo struct {
	clock func() time.Time
	stmt  sq.StatementBuilderType
}

func (p *PhotoRepo) List(ctx context.Context, db *sqlx.DB, user User, paginator database.Paginator) ([]Photo, database.Paginator, error) {
	stmt := p.stmt.
		Select("photos.*").
		From("photos").
		Join("collections c on (photos.collection_id = c.id)").
		Join("users_collections uc on (uc.collection_id = c.id)").
		Where(sq.Eq{"uc.user_id": user.ID})

	sql, args, err := paginator.Paginate(stmt).ToSql()
	if err != nil {
		return nil, paginator, errors.Wrap(err, "could not build query")
	}

	var photos []Photo
	err = db.SelectContext(ctx, &photos, sql, args...)
	if err != nil {
		return nil, paginator, errors.Wrap(err, "could not select rows")
	}

	return photos, paginator, nil
}

// AddPhoto creates a new photo, original rendition, and if applicable, exif records from the given
// reader. Returns the photo instance, the original rendition, or an error.
func (p *PhotoRepo) AddPhoto(ctx context.Context, tx sqlx.ExtContext, storage storage.Backend, collection Collection, upload PhotoUpload) (Photo, Rendition, error) {
	var takenAt *time.Time
	e, err := exif.Decode(upload.Reader)
	if err != nil && exif.IsCriticalError(err) {
		log.Printf("error getting exif tags: %v", err)
	} else {
		if dateTime, err := e.DateTime(); err != nil {
			log.Printf("error getting exif datetime tags: %v", err)
		} else {
			takenAt = &dateTime
		}
	}

	if _, err := upload.Reader.Seek(0, io.SeekStart); err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not rewind")
	}

	photo := Photo{
		Timestamps:     db.JustCreated(p.clock),
		CollectionID:   collection.ID,
		RenditionCount: 1,
		Description:    "",
		Filename:       upload.Filename,
		TakenAt:        takenAt,
		Published:      false,
	}

	photo, err = p.Create(ctx, tx, photo)
	if err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not insert photo")
	}

	renditionConfig, err := FindOriginalRenditionConfiguration(ctx, tx)
	if err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not find rendition config for original")
	}

	rawJpeg, err := jpeg.Decode(upload.Reader)
	if err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not decode jpeg")
	}
	width, height := uint(rawJpeg.Bounds().Dx()), uint(rawJpeg.Bounds().Dy())
	rendition := Rendition{
		Format:                   upload.ContentType,
		Height:                   height,
		Original:                 true,
		PhotoID:                  photo.ID,
		RenditionConfigurationID: renditionConfig.ID,
		Timestamps:               db.JustCreated(p.clock),
		Width:                    width,
	}
	rendition, err = InsertRendition(ctx, tx, rendition)
	if err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not insert rendition")
	}

	if _, err := upload.Reader.Seek(0, io.SeekStart); err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not rewind")
	}

	buf, err := ioutil.ReadAll(upload.Reader)
	if err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not read all bytes")
	}

	if err := storage.Store(rendition.ID, buf); err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not store rendition")
	}

	return photo, rendition, nil
}

// Create stores a new photo in the database.
func (p *PhotoRepo) Create(ctx context.Context, tx sqlx.ExtContext, photo Photo) (Photo, error) {
	sql, args, err := p.stmt.Insert("photos").
		Columns("updated_at", "created_at", "collection_id", "rendition_count", "description", "filename", "taken_at", "published").
		Values(photo.UpdatedAt, photo.CreatedAt, photo.CollectionID, photo.RenditionCount, photo.Description, photo.Filename, photo.TakenAt, photo.Published).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return photo, errors.Wrap(err, "could not build query")
	}

	rows, err := tx.QueryxContext(ctx, sql, args...)
	defer rows.Close()
	if err != nil {
		return photo, errors.Wrap(err, "could not insert")
	}

	if !rows.Next() {
		return photo, errors.Wrap(err, "no id returned")
	}

	err = rows.Scan(&photo.ID)
	if err != nil {
		return photo, errors.Wrap(err, "could scan id")
	}

	return photo, nil
}
