package model

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExifTagsFromPhoto(t *testing.T) {
	// TODO: add a simple jpeg to the repo
	data, err := ioutil.ReadFile("../test/integration/files/100x75-with-exif.jpg")
	assert.Nil(t, err)
	tags, err := ExifTagsFromPhoto(data)
	assert.Nil(t, err)

	assert.NotNil(t, tags)
}
