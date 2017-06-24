package model

import (
	"fmt"
	"log"
	"math"

	"github.com/ilikeorangutans/phts/db"
)

type Renditions []Rendition

func (r Renditions) NotEmpty() bool {
	return len(r) > 0
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

type Photo struct {
	db.PhotoRecord
	Renditions Renditions
}

type Rendition struct {
	db.RenditionRecord
}
