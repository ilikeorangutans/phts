package model

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"time"

	"github.com/disintegration/imaging"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/metadata"
	"github.com/jmoiron/sqlx"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"

	sq "github.com/Masterminds/squirrel"
)

// RenditionConfiguration describes how to produce a rendition
type RenditionConfiguration struct {
	db.Record
	db.Timestamps

	CollectionID *int64 `db:"collection_id" json:"collectionID"`
	Height       int    `db:"height" json:"height"`
	Name         string `db:"name" json:"name"`
	Original     bool   `db:"original" json:"original"`
	Private      bool   `db:"private" json:"private"` // TODO rename this to "system" or "reserved" or "locked"
	Quality      int    `db:"quality" json:"quality"`
	Resize       bool   `db:"resize" json:"resize"`
	Width        int    `db:"width" json:"width"`
}

func (r RenditionConfiguration) Process(ctx context.Context, data []byte) (Rendition, []byte, error) {
	orientation := metadata.Horizontal

	reader := bytes.NewReader(data)
	e, err := exif.Decode(reader)
	if err != nil && exif.IsCriticalError(err) {
	} else {
		if orientationTag, err := e.Get(exif.Orientation); err == nil {
			if orientationValue, err := orientationTag.Int(0); err == nil {
				orientation = metadata.ExifOrientation(orientationValue)
			}
		}
	}

	var rendition Rendition
	// TODO move most of this into the rendition
	rawJpeg, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return rendition, nil, errors.Wrap(err, "error decoding jpeg")
	}

	rawJpeg = rotate(rawJpeg, orientation.Angle())

	width, height := uint(rawJpeg.Bounds().Dx()), uint(rawJpeg.Bounds().Dy())
	if orientation.Angle()%180 != 0 {
		width, height = height, width
	}

	binary := data

	if r.Resize {
		// TODO instead of reading from rawJpeg we should take the previous result (which should be smaller than the original, but bigger than this version
		resized := resize.Resize(uint(r.Width), 0, rawJpeg, resize.Lanczos3)
		var b = &bytes.Buffer{}
		if err := jpeg.Encode(b, resized, &jpeg.Options{Quality: r.Quality}); err != nil {
			return rendition, nil, errors.Wrap(err, "could not encode jpeg")
		}
		width = uint(resized.Bounds().Dx())
		height = uint(resized.Bounds().Dy())
		binary = b.Bytes()
	}

	rendition = Rendition{
		Timestamps:               db.JustCreated(time.Now),
		Width:                    width,
		Height:                   height,
		Format:                   "image/jpeg",
		Original:                 false,
		RenditionConfigurationID: r.ID,
	}

	return rendition, binary, nil
}

// FindOriginalRenditionConfiguration finds the single rendition configuration for original renditions.
func FindOriginalRenditionConfiguration(ctx context.Context, dbx sqlx.ExtContext) (RenditionConfiguration, error) {
	var config RenditionConfiguration
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("rendition_configurations").
		Where(sq.Eq{"original": true, "collection_id": nil}).
		Limit(1).
		ToSql()
	if err != nil {
		return config, errors.Wrap(err, "could not create query")
	}

	err = dbx.QueryRowxContext(ctx, sql, args...).StructScan(&config)
	if err != nil {
		return config, errors.Wrap(err, "could not query row")
	}

	return config, nil
}

// FindMissingRenditionConfigurations finds all rendition configurations for which the given photo doesn't have
// a rendition.
func FindMissingRenditionConfigurations(ctx context.Context, dbx sqlx.QueryerContext, photo Photo) ([]RenditionConfiguration, error) {
	sql := `
	  select
		*
	  from
		rendition_configurations
	  where
		(
		  collection_id is null
		  or
		  collection_id = $1
		)
		and
		id not in (
		  select
			rendition_configuration_id
		  from
			renditions
		  where
			photo_id = $2
		)
	`

	rows, err := dbx.QueryxContext(ctx, sql, photo.CollectionID, photo.ID)
	if err != nil {
		return nil, errors.Wrap(err, "could not query")
	}

	var configs []RenditionConfiguration
	for rows.Next() {
		var config RenditionConfiguration
		rows.StructScan(&config)
		configs = append(configs, config)
	}

	return configs, nil
}

// FindApplicableRenditionConfigurations finds all RenditionConfigurations that are applicable to the given
// collection, ignoring the original rendition.
func FindApplicableRenditionConfigurations(ctx context.Context, dbx sqlx.QueryerContext, collection Collection) ([]RenditionConfiguration, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("rendition_configurations").
		Where(
			sq.And{
				sq.Eq{"original": false},
				sq.Or{
					sq.Eq{
						"collection_id": collection.ID,
					},
					sq.Eq{
						"collection_id": nil,
					},
				},
			},
		).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not create query")
	}

	rows, err := dbx.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "could not query")
	}

	var configs []RenditionConfiguration
	for rows.Next() {
		var config RenditionConfiguration
		rows.StructScan(&config)
		configs = append(configs, config)
	}

	return configs, nil
}

func rotate(img image.Image, angle int) image.Image {
	//var result *image.NRGBA
	var result image.Image = img
	switch angle {
	case -90:
		// Angles are opposite as imaging uses counter clockwise angles and we use clockwise.
		result = imaging.Rotate270(img)
	case 90:
		result = imaging.Rotate270(img)
	case 180:
		result = imaging.Rotate180(img)
	default:
	}
	return result
}
