package model

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"

	"github.com/disintegration/imaging"
	"github.com/ilikeorangutans/phts/db"
	"github.com/nfnt/resize"
)

type RenditionConfiguration struct {
	db.RenditionConfigurationRecord
}

type RenditionConfigurations []RenditionConfiguration

func (r RenditionConfigurations) ByID(id int64) (RenditionConfiguration, error) {
	for _, config := range r {
		if config.ID == id {
			return config, nil
		}
	}
	return RenditionConfiguration{}, fmt.Errorf("no rendition configuration with ID %d", id)
}

// Without returns a new set of configurations without the specified excludes.
func (r RenditionConfigurations) Without(exclude RenditionConfigurations) RenditionConfigurations {
	if exclude == nil || len(exclude) == 0 {
		return r
	}

	var result RenditionConfigurations
	for _, a := range r {
		found := false
		for _, b := range exclude {
			if a.ID == b.ID {
				found = true
				break
			}
		}

		if !found {
			result = append(result, a)
		}
	}

	return result
}

func (r RenditionConfigurations) Process(filename string, data []byte, orientation ExifOrientation) (Renditions, error) {
	// TODO sort configs by size: big -> small
	var renditions Renditions
	for _, config := range r {
		rawJpeg, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}

		log.Printf("adding %s, orientation: %s", filename, orientation)
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
				return nil, err
			}
			width = uint(resized.Bounds().Dx())
			height = uint(resized.Bounds().Dy())
			binary = b.Bytes()
		}

		record := db.RenditionRecord{
			Original: !config.Resize,
			Width:    width,
			Height:   height,
			Format:   "image/jpeg",
			RenditionConfigurationID: config.ID,
		}

		renditions = append(renditions, Rendition{record, binary})
	}
	return renditions, nil
}

type RenditionConfigurationsBySizeDescending struct {
	RenditionConfigurations
}

func (r RenditionConfigurationsBySizeDescending) Len() int { return len(r.RenditionConfigurations) }
func (r RenditionConfigurationsBySizeDescending) Swap(i, j int) {
	r.RenditionConfigurations[i], r.RenditionConfigurations[j] = r.RenditionConfigurations[j], r.RenditionConfigurations[i]
}
func (r RenditionConfigurationsBySizeDescending) Less(i, j int) bool {
	return r.RenditionConfigurations[i].Area() > r.RenditionConfigurations[j].Area()
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
