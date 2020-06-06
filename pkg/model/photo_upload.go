package model

import (
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
)

var ErrInvalidFiletype = errors.New("invalid file type")

func FromReader(reader io.ReadSeeker, filename string) (PhotoUpload, error) {
	mime, err := mimetype.DetectReader(reader)
	if err != nil {
		return PhotoUpload{}, errors.Wrap(err, "could not read for mime type detection")
	}
	if !strings.HasPrefix(mime.String(), "image/") {
		return PhotoUpload{}, ErrInvalidFiletype
	}

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return PhotoUpload{}, errors.Wrap(err, "error rewinding")
	}

	return PhotoUpload{
		Filename:    filename,
		Reader:      reader,
		ContentType: mime.String(),
	}, nil
}

type PhotoUpload struct {
	Filename    string
	Reader      io.ReadSeeker
	ContentType string
}
