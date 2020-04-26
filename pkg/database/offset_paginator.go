package database

import (
	"fmt"
	"html/template"
	"math"
	"net/url"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/ilikeorangutans/phts/pkg/util"
)

// OffsetPaginatorOpts describes the pagination options for a collection. It describes defaults, min, and max values for limit values and what
// columns are valid for ordering.
type OffsetPaginatorOpts struct {
	MinLimit, DefaultLimit, MaxLimit int
	ValidOrderColumns                []string
	DefaultOrderColumn               string
	DefaultOrder                     string
}

// PaginatorFromQuery creates a new OffsetPaginator using the defaults in this opts and the query values.
func (o OffsetPaginatorOpts) PaginatorFromQuery(query url.Values) OffsetPaginator {
	// TODO pull out order asc|desc
	return OffsetPaginator{
		Offset:  uint64(o.Offset(query.Get("offset"))),
		Limit:   uint64(o.Limit(query.Get("limit"))),
		OrderBy: o.OrderByColumn(query.Get("orderBy")),
		Order:   o.Order(query.Get("order")),
	}
}

func (o OffsetPaginatorOpts) Offset(desired string) int {
	offset, err := strconv.Atoi(desired)
	if err != nil || offset < 0 {
		return 0
	}

	return offset
}

func (o OffsetPaginatorOpts) Limit(desired string) int {
	limit, err := strconv.Atoi(desired)
	if err != nil {
		return o.DefaultLimit
	}

	if limit < o.MinLimit || limit > o.MaxLimit {
		return o.DefaultLimit
	}

	return limit
}

func (o OffsetPaginatorOpts) OrderByColumn(desired string) string {
	if util.StringSliceContains(o.ValidOrderColumns, desired) {
		return desired
	}

	return o.DefaultOrderColumn
}

func (o OffsetPaginatorOpts) Order(desired string) string {
	if desired == "asc" || desired == "desc" {
		return desired
	}

	return o.DefaultOrder
}

type OffsetPaginator struct {
	Offset  uint64
	Limit   uint64
	OrderBy string
	Order   string
	Count   uint64
}

func (o OffsetPaginator) WithCount(count uint64) OffsetPaginator {
	// TODO check if offset/limit make sense with the given count and clamp if necessary
	o.Count = count
	return o
}

// Page returns the current page number in the pagination sequence
func (o OffsetPaginator) Page() uint {
	if o.Count == 0 {
		return 1
	}

	return uint(math.Ceil(float64(o.Offset)/float64(o.Limit))) + 1
}

// PageCount returns the maximum number of pages
func (o OffsetPaginator) PageCount() uint {
	if o.Count == 0 {
		return 1
	}

	return uint(math.Ceil(float64(o.Count) / float64(o.Limit)))
}

func (o OffsetPaginator) HasPrev() bool {
	return o.Offset > 0
}

func (o OffsetPaginator) HasNext() bool {
	if o.Count > 0 {
		return o.Offset+o.Limit < o.Count
	} else {
		// This is a guess
		return true
	}
}

// Next returns a new paginator for the next page
func (o OffsetPaginator) Next() OffsetPaginator {
	next := o
	next.Offset = next.Offset + next.Limit
	return next
}

// Prev returns a new paginator for the previous page
func (o OffsetPaginator) Prev() OffsetPaginator {
	prev := o
	prev.Offset = 0
	if prev.Offset > prev.Limit {
		prev.Offset = prev.Offset - prev.Limit
	}
	return prev
}

// Paginate takes a sql builder and adds the sort and order conditions describe by this paginator.
func (o OffsetPaginator) Paginate(input sq.SelectBuilder) sq.SelectBuilder {
	return input.OrderBy(fmt.Sprintf("%s %s", o.OrderBy, o.Order)).Limit(o.Limit).Offset(o.Offset)
}

func (o OffsetPaginator) QueryString() template.URL {
	v := url.Values{}
	v.Set("limit", strconv.FormatUint(o.Limit, 10))
	v.Set("offset", strconv.FormatUint(o.Offset, 10))
	v.Set("orderBy", o.OrderBy)
	v.Set("order", o.Order)
	return template.URL(v.Encode())
}

func (o OffsetPaginator) String() string {
	return fmt.Sprintf("OffsetPaginator{limit=%d,offset=%d,orderBy=%s,count=%d}", o.Limit, o.Offset, o.OrderBy, o.Count)
}
