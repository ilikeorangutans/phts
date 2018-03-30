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

func (p Paginator) PaginateSqurrel(input sq.SelectBuilder) sq.SelectBuilder {
	query := input

	if p.PrevTimestamp != nil {
		// TODO sadly we can't write composite values with squirrel, so we have to emulate
		query = query.
			Where(
				sq.LtOrEq{
					p.prefixedColumn(p.Column): p.PrevTimestamp,
				},
			).
			Where(
				sq.Or{
					sq.Lt{
						p.prefixedColumn(p.Column): p.PrevTimestamp,
					},
					sq.And{
						sq.Eq{
							p.prefixedColumn(p.Column): p.PrevTimestamp,
						},
						sq.Lt{
							p.prefixedColumn("id"): p.PrevID,
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

var regex = regexp.MustCompile("\\$[0-9]+")

func (p Paginator) Paginate(query string, args ...interface{}) (string, []interface{}) {
	next := findNextPlaceholder(query)
	var buffer bytes.Buffer
	buffer.WriteString(query)

	if p.PrevTimestamp != nil {
		buffer.WriteString(" AND (")
		buffer.WriteString(p.prefixedColumn(p.Column))
		buffer.WriteString(",")
		buffer.WriteString(p.prefixedColumn("id"))
		buffer.WriteString(")")
		buffer.WriteString(p.Direction.AfterRelation())
		buffer.WriteString("(")
		buffer.WriteString("$")
		buffer.WriteString(strconv.Itoa(next))
		next++
		buffer.WriteString(",$")
		buffer.WriteString(strconv.Itoa(next))
		next++
		buffer.WriteString(")")

		args = append(args, *p.PrevTimestamp)
		args = append(args, p.PrevID)

	}

	buffer.WriteString(" ORDER BY ")
	buffer.WriteString(p.prefixedColumn(p.Column))
	buffer.WriteString(" ")
	buffer.WriteString(string(p.Direction))
	buffer.WriteString(",")
	buffer.WriteString(p.prefixedColumn("id"))
	buffer.WriteString(" ")
	buffer.WriteString(string(p.Direction))
	buffer.WriteString(" LIMIT $")
	buffer.WriteString(strconv.Itoa(next))

	return buffer.String(), append(args, p.Count)
}

func (p Paginator) prefixedColumn(name string) string {
	if len(p.ColumnPrefix) == 0 {
		return name
	}

	return fmt.Sprintf("%s.%s", p.ColumnPrefix, name)
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
