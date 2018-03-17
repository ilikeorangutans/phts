package storage

import (
	"bytes"
	"context"
	"io"
	"log"
	"strconv"

	"cloud.google.com/go/storage"
)

func NewGCSBackend(projectID string, ctx context.Context, bucketName string) (Backend, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)

	if err := bucket.Create(ctx, projectID, nil); err != nil {
		log.Fatal(err)
	}

	log.Printf("Google cloud storage backend ready in project %s in bucket %s", projectID, bucketName)

	return &GCSBackend{
		ctx:    ctx,
		client: client,
		bucket: bucket,
	}, nil
}

type GCSBackend struct {
	ctx       context.Context
	projectID string
	client    *storage.Client
	bucket    *storage.BucketHandle
}

func (b *GCSBackend) Store(id int64, data []byte) error {
	name := strconv.FormatInt(id, 10)
	obj := b.bucket.Object(name)
	writer := obj.NewWriter(b.ctx)
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func (b *GCSBackend) Get(id int64) ([]byte, error) {
	obj := b.bucket.Object(strconv.FormatInt(id, 10))
	reader, err := obj.NewReader(b.ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	buffer := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buffer, reader)
	return buffer.Bytes(), err
}

func (b *GCSBackend) Delete(id int64) error {
	obj := b.bucket.Object(strconv.FormatInt(id, 10))
	return obj.Delete(b.ctx)
}
