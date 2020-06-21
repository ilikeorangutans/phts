package server

import (
	"context"
	"time"

	"github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//  StartRenditionUpdateQueueHandler starts a go routine to continuously check for missing renditions and numWorkers go routines to process rendition updates.
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
	l := r.logger.With().Int64("photo-id", req.Photo.ID).Logger()
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

	for _, config := range missingRenditions {
		l.Debug().Str("rendition", config.Name).Msg("generating rendition")
		err := r.processRenditionUpdate(ctx, config, req.Photo, data)
		if err != nil {
			l.Warn().Err(err).Int64("rendition-configuration-id", config.ID).Msg("could not process config")
			continue
		}
		l.Debug().Str("rendition", config.Name).Msg("rendition created")
	}
	l.Debug().Msg("renditions up to date")

	return nil
}

func (r *renditionUpdateWorker) processRenditionUpdate(ctx context.Context, config model.RenditionConfiguration, photo model.Photo, data []byte) error {
	photoRepo := model.NewPhotoRepo()
	rendition, binary, err := config.Process(ctx, data)
	if err != nil {
		return errors.Wrap(err, "could not process config")
	}

	tx, err := r.dbx.Beginx()
	if err != nil {
		return errors.Wrap(err, "could not start transaction")
	}

	_, rendition, err = photoRepo.AddRendition(ctx, tx, photo, rendition)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "error adding rendition")
	}

	err = r.backend.Store(rendition.ID, binary)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "could not store binary")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		r.backend.Delete(rendition.ID)
		return errors.Wrap(err, "could not commit")
	}

	return nil
}
