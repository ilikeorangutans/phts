package db

import (
	"fmt"
	"log"
)

type Direction string

func (d Direction) AfterRelation() string {
	switch d {
	case Asc:
		return ">"
	case Desc:
		return "<"
	default:
		log.Panicf("Invalid sort direction %q", d)
		return ""
	}
}

func (d Direction) AddToColumn(column string) string {
	return fmt.Sprintf("%s %s", column, d)
}

const (
	Asc  Direction = "ASC"
	Desc           = "DESC"
)
