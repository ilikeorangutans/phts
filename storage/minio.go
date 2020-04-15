package storage

import (
	"bytes"
	"io"
	"log"
	"strconv"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
)

func NewMinIOBackend(endpoint, accessKey, secretKey, bucket string, useSSL bool) (Backend, error) {
	log.Printf("starting minio backend with endpoint %s using bucket %s", endpoint, bucket)
	client, err := minio.New(endpoint, accessKey, secretKey, useSSL)
	if err != nil {
		return nil, errors.Wrap(err, "could not create minio backend")
	}

	bucketExists, err := client.BucketExists(bucket)
	if err != nil {
		return nil, errors.Wrap(err, "could not check if bucket exists")
	}
	if !bucketExists {
		log.Printf("creating bucket %s", bucket)
		if err := client.MakeBucket(bucket, ""); err != nil {
			return nil, errors.Wrap(err, "could not create bucket")
		}
	}

	return &MinIOBackend{
		client: client,
		bucket: bucket,
	}, nil
}

type MinIOBackend struct {
	client *minio.Client
	bucket string
}

func (m *MinIOBackend) Store(id int64, data []byte) error {
	_, err := m.client.PutObject(m.bucket, strconv.FormatInt(id, 10), bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "could not store object")
	}
	return nil
}

func (m *MinIOBackend) Get(id int64) ([]byte, error) {
	obj, err := m.client.GetObject(m.bucket, strconv.FormatInt(id, 10), minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "could not get object")
	}
	buffer := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buffer, obj)
	if err != nil {
		return nil, errors.Wrap(err, "could not copy bytes")
	}
	return buffer.Bytes(), nil
}

func (m *MinIOBackend) Delete(id int64) error {
	err := m.client.RemoveObject(m.bucket, strconv.FormatInt(id, 10))
	if err != nil {
		return errors.Wrap(err, "could not delete object")
	}
	return nil
}
