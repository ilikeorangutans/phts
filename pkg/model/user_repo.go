package model

import (
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

var UsersPaginator = database.OffsetPaginatorOpts{
	MinLimit:           1,
	DefaultLimit:       10,
	MaxLimit:           100,
	ValidOrderColumns:  []string{"id", "created_at", "updated_at", "email"},
	DefaultOrderColumn: "id",
	DefaultOrder:       "asc",
}

func NewUserRepo(db db.DB) *UserRepo {
	return &UserRepo{
		db:    db,
		clock: time.Now,
	}
}

type UserRepo struct {
	db    db.DB
	clock func() time.Time
}

func (u *UserRepo) List(paginator database.OffsetPaginator) ([]User, database.OffsetPaginator, error) {
	var users []User
	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	var count uint64
	err := stmt.RunWith(u.db).Select("count(*)").From("users").QueryRow().Scan(&count)
	if err != nil {
		return users, paginator, errors.Wrap(err, "could not list users")
	}
	paginator = paginator.WithCount(count)

	query := stmt.Select("*").From("users")
	query = paginator.Paginate(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return users, paginator, errors.Wrap(err, "could not list users")
	}

	err = u.db.Select(&users, sql, args...)
	if err != nil {
		return users, paginator, errors.Wrap(err, "could not list users")
	}

	return users, paginator, nil
}
