package model

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/ilikeorangutans/phts/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Share struct {
	db.Record
	db.Timestamps
	PhotoID      int64  `db:"photo_id" json:"photoID"`
	CollectionID int64  `db:"collection_id" json:"collectionID"`
	ShareSiteID  int64  `db:"share_site_id" json:"shareSiteID"`
	Slug         string `db:"slug" json:"slug"`
}

func FindSharedPhotoBySlug(ctx context.Context, tx sqlx.QueryerContext, shareSite ShareSite, slug string) (ShareWithPhotos, error) {
	var shareWithPhotos ShareWithPhotos
	share, err := FindShareBySiteAndSlug(ctx, tx, shareSite, slug)
	if err != nil {
		return shareWithPhotos, errors.Wrap(err, "could not find share for slug")
	}

	renditionConfigs, err := FindRenditionConfigurationsForShare(ctx, tx, share)
	if err != nil {
		return shareWithPhotos, errors.Wrap(err, "could not find rendition configurations for share")
	}

	// TODO right now we only support single photos in a share
	photo, err := NewPhotoRepo().FindByID(ctx, tx, share.PhotoID)
	renditions, err := FindRenditionsForPhoto(ctx, tx, photo, renditionConfigs...)

	shareWithPhotos.Share = share
	shareWithPhotos.RenditionConfigurations = renditionConfigs
	shareWithPhotos.Photos = []PhotoWithRenditions{
		{
			Photo:      photo,
			Renditions: renditions,
		},
	}

	return shareWithPhotos, nil
}

// ShareWithPhotos holds a share, its photos, and the associated rendition configurations.
type ShareWithPhotos struct {
	Share                   Share
	Photos                  []PhotoWithRenditions
	RenditionConfigurations []RenditionConfiguration
}

type PhotoWithRenditions struct {
	Photo      Photo
	Renditions []Rendition
}

func FindRenditionConfigurationsForShare(ctx context.Context, tx sqlx.QueryerContext, share Share) ([]RenditionConfiguration, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("rc.*").
		From("rendition_configurations AS rc").
		Join("share_rendition_configurations AS src ON rc.id = src.rendition_configuration_id").
		Where(sq.Eq{"src.share_id": share.ID}).
		Limit(10). // TODO rather arbitrary, but do we really need more than ten?
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "could not create query")
	}

	rows, err := tx.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "could not query rows")
	}
	defer rows.Close()

	var configs []RenditionConfiguration
	for rows.Next() {
		var config RenditionConfiguration
		err := rows.StructScan(&config)
		if err != nil {
			return nil, errors.Wrap(err, "could not query rows")
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// FindShareBySiteAndSlug looks up a share by the given share site and slug.
func FindShareBySiteAndSlug(ctx context.Context, tx sqlx.QueryerContext, shareSite ShareSite, slug string) (Share, error) {
	var share Share
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From("shares").
		Where(sq.Eq{
			"share_site_id": shareSite.ID,
			"slug":          slug,
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return share, errors.Wrap(err, "could not build query")
	}

	log.Printf("share site %v", shareSite)
	log.Printf("slug %s", slug)
	log.Printf("sql %s", sql)

	err = tx.QueryRowxContext(ctx, sql, args...).StructScan(&share)
	if err != nil {
		return share, errors.Wrap(err, "could not query row")
	}

	return share, nil
}
