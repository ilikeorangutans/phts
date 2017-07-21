package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func paginator() Paginator {
	var p = Paginator{
		Direction: Desc,
		Column:    "updated_at",
		Count:     10,
	}
	return p
}

func TestPaginatorFirstPage(t *testing.T) {
	p := paginator()
	query, fields := p.Paginate("SELECT * FROM foo WHERE blargh=$1", 1)

	assert.Equal(t, "SELECT * FROM foo WHERE blargh=$1 ORDER BY updated_at DESC, id DESC LIMIT $2", query)
	assert.Equal(t, []interface{}{1, uint(10)}, fields)
}

func TestPaginatorWithPrevValues(t *testing.T) {
	p := paginator()
	ts := time.Now()
	p.PrevTimestamp = &ts
	p.PrevID = 17
	query, fields := p.Paginate("SELECT * FROM foo WHERE blargh=$1", 1)

	assert.Equal(t, "SELECT * FROM foo WHERE blargh=$1 AND (updated_at,id)<($2,$3) ORDER BY updated_at DESC, id DESC LIMIT $4", query)
	expected := []interface{}{1, ts, int64(17), uint(10)}
	assert.Equal(t, expected, fields)
}

func TestPaginatorMultipleFieldWithoutPrevValues(t *testing.T) {
	p := paginator()
	query, fields := p.Paginate("SELECT * FROM foo WHERE blargh=$1 AND blurgh = $2", 1, 2)

	assert.Equal(t, "SELECT * FROM foo WHERE blargh=$1 AND blurgh = $2 ORDER BY updated_at DESC, id DESC LIMIT $3", query)
	assert.Equal(t, []interface{}{1, 2, uint(10)}, fields)
}

func TestPaginatorMultipleField(t *testing.T) {
	p := paginator()
	ts := time.Now()
	p.PrevTimestamp = &ts
	p.PrevID = 17
	query, fields := p.Paginate("SELECT * FROM foo WHERE blargh=$1 AND blurgh = $2", 1, 2)

	assert.Equal(t, "SELECT * FROM foo WHERE blargh=$1 AND blurgh = $2 AND (updated_at,id)<($3,$4) ORDER BY updated_at DESC, id DESC LIMIT $5", query)
	assert.Equal(t, []interface{}{1, 2, ts, int64(17), uint(10)}, fields)
}
