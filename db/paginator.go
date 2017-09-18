package db

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"
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

const (
	Asc  Direction = "ASC"
	Desc           = "DESC"
)

func PaginatorFromRequest(query url.Values) Paginator {
	p := NewPaginator()
	var err error
	if p.PrevID, err = strconv.ParseInt(query.Get("prevID"), 10, 64); err != nil {
		return p
	}

	if timestampString := query.Get("prevTimestamp"); len(timestampString) > 0 {
		t, err := time.Parse(time.RFC3339, timestampString)
		if err != nil {
			return p
		}
		p.PrevTimestamp = &t
		return p
	} else {
		p.PrevValue = query.Get("prevValue")
	}

	return p
}

func NewPaginator() Paginator {
	return Paginator{
		Direction: Desc,
		Column:    "updated_at",
		Count:     10,
	}
}

type Paginator struct {
	Direction     Direction   `json:"direction"`
	Column        string      `json:"column"`
	Count         uint        `json:"count"`
	PrevValue     interface{} `json:"prev_value"`
	PrevTimestamp *time.Time  `json:"prev_timestamp"`
	PrevID        int64       `json:"prev_id"`
}

func (p Paginator) QueryString() template.URL {
	prevField := ""
	prevValue := ""
	if p.PrevTimestamp != nil {
		prevField = "prevTimestamp"
		prevValue = p.PrevTimestamp.Format(time.RFC3339)
	}
	return template.URL(fmt.Sprintf("prevID=%d&%s=%s", p.PrevID, prevField, prevValue))
}

var regex = regexp.MustCompile("\\$[0-9]+")

func (p Paginator) Paginate(query string, args ...interface{}) (string, []interface{}) {
	next := findNextPlaceholder(query)
	var buffer bytes.Buffer
	buffer.WriteString(query)

	if p.PrevTimestamp != nil {
		buffer.WriteString(" AND (")
		buffer.WriteString(p.Column)
		buffer.WriteString(",id)")
		buffer.WriteString(p.Direction.AfterRelation())
		buffer.WriteString("(")
		buffer.WriteString("$")
		buffer.WriteString(strconv.Itoa(next))
		next += 1
		buffer.WriteString(",$")
		buffer.WriteString(strconv.Itoa(next))
		next += 1
		buffer.WriteString(")")

		args = append(args, *p.PrevTimestamp)
		args = append(args, p.PrevID)

	}

	buffer.WriteString(" ORDER BY ")
	buffer.WriteString(p.Column)
	buffer.WriteString(" ")
	buffer.WriteString(string(p.Direction))
	buffer.WriteString(", id ")
	buffer.WriteString(string(p.Direction))
	buffer.WriteString(" LIMIT $")
	buffer.WriteString(strconv.Itoa(next))

	return buffer.String(), append(args, p.Count)
}

func findNextPlaceholder(query string) int {
	vars := regex.FindAllString(query, -1)
	max := 0
	for _, x := range vars {
		cur, err := strconv.Atoi(x[1:])
		if err != nil {
			log.Panic(err)
		}
		if cur > max {
			max = cur
		}
	}

	return max + 1
}
