package model

import (
	"fmt"
	"log"
	"math"

	"github.com/ilikeorangutans/phts/db"
)

type Renditions []Rendition

type Rendition struct {
	db.RenditionRecord
	data []byte // Not meant to actually contain the image binary, just for passing it around
}

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
