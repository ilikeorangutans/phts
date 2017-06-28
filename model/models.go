package model

import (
	"fmt"
	"log"
	"math"

	"github.com/ilikeorangutans/phts/db"
	"github.com/rwcarlsen/goexif/tiff"
)

type Renditions []Rendition

func (r Renditions) NotEmpty() bool {
	return len(r) > 0
}

func (r Renditions) Empty() bool {
	return len(r) == 0
}

func (r Renditions) Smallest() Rendition {
	if len(r) == 0 {
		log.Fatal(fmt.Errorf("Cannot call Smallest() on empty set of renditions!"))
	}
	min := uint(math.MaxUint32)
	index := 0
	for i, rendition := range r {
		if rendition.Width < min {
			min = rendition.Width
			index = i
		}
	}

	return r[index]
}

func (r Renditions) Large() Rendition {
	if r.Empty() {
		log.Panic("Cannot return large rendition on empty set")
	}

	max := uint(0)
	index := 0
	for i, rendition := range r {
		if !rendition.Original && rendition.Width > max {
			max = rendition.Width
			index = i
		}
	}

	return r[index]
}

type Photo struct {
	db.PhotoRecord
	Renditions Renditions
	Exif       []ExifTag
}

type Rendition struct {
	db.RenditionRecord
}

type ExifTag struct {
	db.ExifRecord
}

func (e ExifTag) String() string {
	switch e.Type {
	case tiff.DTAscii:
		return fmt.Sprintf("%s", e.StringValue)
	case tiff.DTLong, tiff.DTShort:
		return fmt.Sprintf("%d", e.Num)
	case tiff.DTRational, tiff.DTSRational:
		return fmt.Sprintf("%d/%d", e.Num, e.Denominator)
	}
	return "unknown"
}
