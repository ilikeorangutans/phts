package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

var queryBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func paginator() Paginator {
	var p = Paginator{
		Direction: Desc,
		Column:    "updated_at",
		Count:     10,
	}
	return p
}

func TestPaginatorDefaultValues(t *testing.T) {
	p := paginator()
	input := queryBuilder.Select("*").From("test")

	result := p.Paginate(input)

	query, args, err := result.ToSql()
	assert.Nil(t, err)
	assert.Equal(t, "SELECT * FROM test ORDER BY updated_at DESC, id DESC LIMIT 10", query)
	var expectedArgs []interface{} = nil
	assert.Equal(t, expectedArgs, args)
}

func TestPaginatorSquirrelWithPrevValues(t *testing.T) {
	p := paginator()
	ts := time.Now()
	p.PrevTimestamp = &ts
	p.PrevID = 17
	input := queryBuilder.Select("*").From("test")
	input = input.Where(sq.Eq{"blargh": 12})

	result := p.Paginate(input)

	query, args, err := result.ToSql()
	assert.Nil(t, err)
	assert.Equal(t, "SELECT * FROM test WHERE blargh = $1 AND updated_at <= $2 AND (updated_at < $3 OR (updated_at = $4 AND id < $5)) ORDER BY updated_at DESC, id DESC LIMIT 10", query)
	var expectedArgs []interface{} = []interface{}{12, p.PrevTimestamp, p.PrevTimestamp, p.PrevTimestamp, int64(17)}
	assert.Equal(t, args, expectedArgs)
}
