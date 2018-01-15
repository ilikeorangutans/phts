package queue

import (
	"context"

	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
)

type ResizeQueue interface {
	Enqueue(userID int64, photo model.Photo, renditionConfiguration model.RenditionConfiguration) error
}

func NewResizeQueue(context context.Context) ResizeQueue {
	return &resizeQueueImpl{}
}

type resizeQueueImpl struct {
	storage storage.Backend
}

func (q *resizeQueueImpl) Enqueue(userID int64, photo model.Photo, renditionConfiguration model.RenditionConfiguration) error {
	return nil
}
