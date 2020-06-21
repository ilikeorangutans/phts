package server

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/metadata"
	"github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/pkg/errors"

	"github.com/disintegration/imaging"
	"github.com/jmoiron/sqlx"
	"github.com/nfnt/resize"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rwcarlsen/goexif/exif"
)

func StartRenditionUpdateQueueHandler(ctx context.Context, dbx *sqlx.DB, backend storage.Backend, queue chan model.RenditionUpdateRequest, numWorkers uint, frequency time.Duration) {
	go enqueueMissingRenditions(ctx, dbx, queue, frequency)

	// TODO make the number of workers configurable
	for i := uint(0); i < numWorkers; i++ {
		worker := newRenditionUpdateWorker(dbx, backend, queue, i)
		go worker.start(ctx)
	}
}

// enqueueMissingRenditions finds photos with missing renditions and enqueues them for processing.
func enqueueMissingRenditions(ctx context.Context, dbx *sqlx.DB, queue chan model.RenditionUpdateRequest, frequency time.Duration) {
	log.Debug().Dur("frequency", frequency).Msg("scanning for missing renditions")
	photoRepo := model.NewPhotoRepo()
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			// TODO this sometimes picks up photos we are already processing
			// It's not bad, but annoying
			photos, err := photoRepo.FindPhotosWithMissingRenditions(ctx, dbx, 20)
			if err != nil {
				log.Warn().Err(err).Msg("error scanning for photos with missing renditions")
				continue
			}

			if len(photos) == 0 {
				continue
			}

			log.Debug().Int("count", len(photos)).Msg("found photos with missing renditions")

			for _, photo := range photos {
				queue <- model.RenditionUpdateRequest{
					Photo: photo,
				}
			}

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// newRenditionUpdateWorker creates a new worker.
func newRenditionUpdateWorker(dbx *sqlx.DB, backend storage.Backend, queue <-chan model.RenditionUpdateRequest, id uint) *renditionUpdateWorker {
	return &renditionUpdateWorker{
		dbx:     dbx,
		backend: backend,
		logger:  log.With().Uint("worker-id", id).Logger(),
		queue:   queue,
	}
}

type renditionUpdateWorker struct {
	logger  zerolog.Logger
	dbx     *sqlx.DB
	backend storage.Backend
	queue   <-chan model.RenditionUpdateRequest
}

// start makes the worker consume events from the queue until the ctx says to stop.
func (r *renditionUpdateWorker) start(ctx context.Context) {
	r.logger.Debug().Msg("worker starting")
	for {
		select {
		case <-ctx.Done():
			r.logger.Debug().Msg("worker shutting down")
			return
		case req := <-r.queue:
			ctx, cancel := context.WithTimeout(r.logger.WithContext(ctx), 60*time.Second)
			defer cancel()

			err := r.processUpdateRequest(ctx, req)
			if err != nil {
				r.logger.Warn().Int64("photo-id", req.Photo.ID).Err(err)
			}
		}
	}
}

func (r *renditionUpdateWorker) processUpdateRequest(ctx context.Context, req model.RenditionUpdateRequest) error {
	l := log.Ctx(ctx).With().Int64("photo-id", req.Photo.ID).Logger()
	l.Debug().Msg("processing queue entry")

	missingRenditions, err := model.FindMissingRenditionConfigurations(ctx, r.dbx, req.Photo)
	if err != nil {
		return errors.Wrap(err, "could not find original rendition")
	}

	if len(missingRenditions) == 0 {
		return nil
	}

	original, err := model.FindOriginalRenditionByPhoto(ctx, r.dbx, req.Photo)
	if err != nil {
		return errors.Wrap(err, "could not find original rendition")
	}

	data, err := r.backend.Get(original.ID)
	if err != nil {
		return errors.Wrap(err, "error fetching original binary")
	}

	orientation := metadata.Horizontal

	reader := bytes.NewReader(data)
	e, err := exif.Decode(reader)
	if err != nil && exif.IsCriticalError(err) {
		l.Debug().Err(err).Msg("error getting exif tags")
	} else {
		if orientationTag, err := e.Get(exif.Orientation); err == nil {
			if orientationValue, err := orientationTag.Int(0); err == nil {
				orientation = metadata.ExifOrientation(orientationValue)
			}
		}
	}

	photoRepo := model.NewPhotoRepo()

	for _, config := range missingRenditions {
		// TODO move most of this into the rendition
		l.Debug().Str("rendition", config.Name).Msg("generating rendition")
		rawJpeg, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return errors.Wrap(err, "error decoding jpeg")
		}

		rawJpeg = rotate(rawJpeg, orientation.Angle())

		width, height := uint(rawJpeg.Bounds().Dx()), uint(rawJpeg.Bounds().Dy())
		if orientation.Angle()%180 != 0 {
			width, height = height, width
		}

		binary := data

		if config.Resize {
			// TODO instead of reading from rawJpeg we should take the previous result (which should be smaller than the original, but bigger than this version
			resized := resize.Resize(uint(config.Width), 0, rawJpeg, resize.Lanczos3)
			var b = &bytes.Buffer{}
			if err := jpeg.Encode(b, resized, &jpeg.Options{Quality: config.Quality}); err != nil {
				l.Warn().Err(err).Msg("error encoding jpeg")
				continue
			}
			width = uint(resized.Bounds().Dx())
			height = uint(resized.Bounds().Dy())
			binary = b.Bytes()
		}

		rendition := model.Rendition{
			Timestamps:               db.JustCreated(time.Now),
			Width:                    width,
			Height:                   height,
			Format:                   "image/jpeg",
			Original:                 false,
			RenditionConfigurationID: config.ID,
		}

		// TODO run in transaction here
		_, rendition, err = photoRepo.AddRendition(ctx, r.dbx, req.Photo, rendition)
		if err != nil {
			l.Warn().Err(err).Int64("rendition-configuration-id", config.ID).Msg("error adding rendition")
			continue
		}

		err = r.backend.Store(rendition.ID, binary)
		if err != nil {
			l.Warn().Err(err).Msg("could not store binary")
			continue
		}

		l.Debug().Str("rendition", config.Name).Msg("rendition created")
	}
	l.Debug().Msg("renditions up to date")

	return nil
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
