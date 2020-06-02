package db

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

type UserDB interface {
	Save(*UserRecord) error
	FindByEmail(string) (*UserRecord, error)
	FindByID(int64) (*UserRecord, error)
}

func NewUserDB(db DB) UserDB {
	return &userSQLDB{
		db:    db,
		clock: time.Now,
		sql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type userSQLDB struct {
	db    DB
	clock Clock
	sql   sq.StatementBuilderType
}

func (u *userSQLDB) Save(record *UserRecord) error {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(u.clock)

		query, args, err := u.sql.
			Update("users").
			Set("password", record.Password).
			Set("updated_at", record.UpdatedAt).Where(sq.Eq{"id": record.ID}).
			ToSql()
		if err != nil {
			return err
		}

		err = checkResult(u.db.Exec(query, args...))
	} else {
		record.Timestamps = JustCreated(u.clock)
		query, args, err := u.sql.
			Insert("users").
			Columns("email", "password", "created_at", "updated_at").
			Values(
				record.Email,
				record.Password,
				record.CreatedAt.UTC(),
				record.UpdatedAt.UTC(),
			).
			Suffix("RETURNING \"id\"").
			ToSql()

		if err != nil {
			return err
		}

		err = u.db.QueryRow(query, args...).Scan(&record.ID)
	}

	return err
}

func (u *userSQLDB) FindByEmail(email string) (*UserRecord, error) {
	query := u.sql.Select("*").From("users").Where(sq.Eq{"email": email}).Limit(1)

	var record UserRecord
	err := queryAndStructScan(u.db, query, &record)

	return &record, err
}

func (u *userSQLDB) FindByID(id int64) (*UserRecord, error) {
	query := u.sql.Select("*").From("users").Where(sq.Eq{"id": id}).Limit(1)

	var record UserRecord
	err := queryAndStructScan(u.db, query, &record)

	return &record, err
}
