package model

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadWithInvalidFileType(t *testing.T) {
	reader := bytes.NewReader([]byte{0, 0, 0})
	pu, err := FromReader(reader, "dummydata.bin")
	assert.True(t, errors.Is(err, ErrInvalidFiletype))

	log.Printf("%v", pu)
}

func TestUploadWithJpeg(t *testing.T) {
	file, err := os.Open("../../test/integration/files/1x1.jpg")
	assert.NoError(t, err)
	pu, err := FromReader(file, "some.jpg")
	assert.NoError(t, err)

	photo, rendition, err := pu.PhotoAndRendition()
	assert.NoError(t, err)

	assert.Equal(t, "image/jpeg", rendition.Format)
	assert.Equal(t, "some.jpg", photo.Filename)
}
