package model

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/ilikeorangutans/phts/pkg/metadata"
	"github.com/pkg/errors"
)

var ErrInvalidFiletype = errors.New("invalid file type")

func FromFormFile(r *http.Request, field string) (PhotoUpload, error) {
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return PhotoUpload{}, errors.Wrap(err, "could not get form file")
	}

	return FromReader(file, fileHeader.Filename)
}

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

func (p PhotoUpload) PhotoAndRendition() (Photo, Rendition, error) {
	var takenAt *time.Time
	tags, err := metadata.ExifTagsFromPhotoReader(p.Reader)
	if err != nil {
		log.Printf("could not decode exif: %v", err)
	} else {
		takenAtFields := []string{"DateTime", "DateTimeOriginal"}
		for _, field := range takenAtFields {
			if tag, err := tags.ByName(field); err == nil {
				takenAt = tag.DateTime
				break
			}
		}
	}

	if _, err := p.Reader.Seek(0, io.SeekStart); err != nil {
		return Photo{}, Rendition{}, errors.Wrap(err, "could not rewind")
	}

	return Photo{
			RenditionCount: 1,
			Filename:       p.Filename, // TODO we should employ a whitelist here
			TakenAt:        takenAt,
		}, Rendition{
			Original: true,
			Format:   p.ContentType,
		}, nil

}

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
