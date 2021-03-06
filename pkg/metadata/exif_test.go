package metadata

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExifTagsFromPhoto(t *testing.T) {
	data, err := ioutil.ReadFile("../../test/integration/files/100x75-with-exif.jpg")
	assert.Nil(t, err)
	tags, err := ExifTagsFromPhoto(data)
	assert.Nil(t, err)

	assert.NotNil(t, tags)

	dateTaken, err := tags.ByName("DateTime")
	assert.Nil(t, err)
	assert.Equal(t, time.Date(2015, time.August, 1, 19, 50, 5, 0, time.UTC), *dateTaken.DateTime)
}

func TestExifTagsFromPhotoWithoutExif(t *testing.T) {
	data, err := ioutil.ReadFile("../../test/integration/files/1x1.jpg")
	assert.NoError(t, err)
	tags, err := ExifTagsFromPhoto(data)
	assert.Error(t, err)

	assert.Empty(t, tags)
}
