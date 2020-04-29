package model

import (
	"log"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

// UsersPaginator is the default paginator settings for users.
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
		stmt:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type UserRepo struct {
	db    db.DB
	clock func() time.Time
	stmt  sq.StatementBuilderType
}

func (u *UserRepo) NewUser(email string) (User, error) {
	user := User{
		Email: email,
	}
	return u.Create(user)
}

func (u *UserRepo) Create(user User) (User, error) {
	err := u.PurgeExpiredPasswordChangeTokens()
	if err != nil {
		return user, errors.Wrap(err, "could not create user")
	}
	user.Timestamps = db.JustCreated(u.clock)

	tx, err := u.db.Beginx()
	if err != nil {
		return user, errors.Wrap(err, "could begin transaction")
	}

	sql, args, err := u.stmt.Insert("users").
		Columns("created_at", "updated_at", "email", "password", "must_change_password").
		Values(user.CreatedAt, user.UpdatedAt, user.Email, "", true).
		Suffix("returning id").
		ToSql()
	if err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "could not build query")
	}

	row := u.db.QueryRow(sql, args...)
	err = row.Scan(&user.ID)
	if err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "could not get id")
	}

	token, err := security.GenerateRandomString(32)
	if err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "could not generate token")
	}
	sql, args, err = u.stmt.Insert("user_password_change_tokens").
		Columns("user_id", "created_at", "token", "invite").
		Values(user.ID, u.clock(), token, true).
		ToSql()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "could not insert token")
	}
	user.PasswordChangeToken = token

	err = tx.Commit()
	if err != nil {
		return user, errors.Wrap(err, "could not commit transaction")
	}
	return user, nil
}

// List lists users from the database according to the given paginator. Returns a list of users, the updated paginator, or an error.
func (u *UserRepo) List(paginator database.OffsetPaginator) ([]User, database.OffsetPaginator, error) {
	var users []User
	var count uint64
	err := u.stmt.RunWith(u.db).
		Select("count(*)").
		From("users").
		QueryRow().
		Scan(&count)
	if err != nil {
		return users, paginator, errors.Wrap(err, "could not list users")
	}
	paginator = paginator.WithCount(count)

	query := u.stmt.Select("*").From("users")
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

// PurgeExpiredPasswordChangeTokens purges all non-invite password reset tokens that are older than one hour.
func (u *UserRepo) PurgeExpiredPasswordChangeTokens() error {
	cutOff := time.Now().AddDate(0, 0, -1)
	sql, args, err := u.stmt.Delete("user_password_change_tokens").
		Where(sq.Lt{"created_at": cutOff}, sq.Eq{"invite": false}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "could not purge expired tokens")
	}

	result, err := u.db.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "could not purge expired tokens")
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "could not purge expired tokens")
	}
	log.Printf("removed %d expired tokens", deleted)
	return nil
}
