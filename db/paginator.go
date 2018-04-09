package db

import (
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"time"

	sq "gopkg.in/Masterminds/squirrel.v1"
)

func PaginatorFromRequest(query url.Values) Paginator {
	p := NewPaginator()
	var err error
	if query.Get("prevID") != "" {
		if p.PrevID, err = strconv.ParseInt(query.Get("prevID"), 10, 64); err != nil {
			return p
		}
	}

	if count, err := strconv.ParseInt(query.Get("count"), 10, 64); err == nil {
		p.Count = uint(count)
	} else {
		p.Count = 10
	}

	if timestampString := query.Get("prevTimestamp"); len(timestampString) > 0 {
		t, err := time.Parse(time.RFC3339, timestampString)
		if err != nil {
			return p
		}
		p.PrevTimestamp = &t
		return p
	}

	p.PrevValue = query.Get("prevValue")

	return p
}

func NewPaginator() Paginator {
	return Paginator{
		Direction: Desc,
		Column:    "updated_at",
		Count:     10,
	}
}

// TODO need an overall count?
type Paginator struct {
	Direction Direction `json:"direction"`
	Column    string    `json:"column"`
	// Count is the number of records we want to fetch per request
	Count         uint        `json:"count"`
	PrevValue     interface{} `json:"prev_value"`
	PrevTimestamp *time.Time  `json:"prev_timestamp"`
	PrevID        int64       `json:"prev_id"`
	ColumnPrefix  string
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

func (p Paginator) Paginate(input sq.SelectBuilder) sq.SelectBuilder {
	query := input

	prefixedPrimary := p.prefixedColumn("id")
	prefixedIncremental := p.prefixedColumn(p.Column)

	if p.PrevTimestamp != nil {
		// TODO sadly we can't write composite values with squirrel, so we have to emulate
		query = query.
			Where(
				sq.LtOrEq{
					prefixedIncremental: p.PrevTimestamp,
				},
			).
			Where(
				sq.Or{
					sq.Lt{
						prefixedIncremental: p.PrevTimestamp,
					},
					sq.And{
						sq.Eq{
							prefixedIncremental: p.PrevTimestamp,
						},
						sq.Lt{
							prefixedPrimary: p.PrevID,
						},
					},
				},
			)
	}

	return query.
		OrderBy(
			p.Direction.AddToColumn(p.prefixedColumn(p.Column)),
			p.Direction.AddToColumn("id"),
		).
		Limit(uint64(p.Count))
}

func (p Paginator) prefixedColumn(name string) string {
	if len(p.ColumnPrefix) == 0 {
		return name
	}

	return fmt.Sprintf("%s.%s", p.ColumnPrefix, name)
}
