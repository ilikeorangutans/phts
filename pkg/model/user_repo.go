package model

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/jmoiron/sqlx"
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

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db:           db,
		clock:        time.Now,
		stmt:         sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		randomString: security.GenerateRandomString,
	}
}

type UserRepo struct {
	db           *sqlx.DB
	clock        func() time.Time
	stmt         sq.StatementBuilderType
	newPassword  func(string) (security.Password, error)
	randomString func(int) (string, error)
}

// FindByEmail finds a user by their email.
func (u *UserRepo) FindByEmail(email string) (User, error) {
	sql, args := u.stmt.Select("*").
		From("users").
		Where(sq.Eq{"email": strings.ToLower(email)}).
		Limit(1).
		MustSql()

	var user User
	err := u.db.Get(&user, sql, args...)
	if err != nil {
		return user, errors.Wrap(err, "could not select user")
	}

	return user, nil
}

// NewUser creates a new user record with the given string, persists it in the database and creates a new invite token.
func (u *UserRepo) NewUser(email string) (User, error) {
	user := User{
		Email: email,
	}
	return u.Create(user)
}

// ActivateInvite activates the user with the given inviteID by setting the specified password and removing the
// invite.
func (u *UserRepo) ActivateInvite(ctx context.Context, inviteID, email, name, password string) (User, error) {
	user, err := u.ByInviteID(ctx, inviteID)
	if err != nil {
		return user, errors.Wrap(err, "cannot find invite")
	}

	user.Email = strings.TrimSpace(email)
	user.Name = strings.TrimSpace(name)
	user.Password, err = security.NewPassword(password)
	if err != nil {
		return user, errors.Wrap(err, "could not update password")
	}

	user.MustChangePassword = false

	log.Printf("activating user %v", user)

	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return user, errors.Wrap(err, "could not begin transaction")
	}

	user, err = u.update(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "updating user failed")
	}

	sql, args, err := u.stmt.Delete("user_password_change_tokens").
		Where(sq.Eq{"user_id": user.ID, "token": inviteID}).
		ToSql()
	if err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "creating query to delete token")
	}

	result, err := u.db.ExecContext(ctx, sql, args...)
	if err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "deleting token")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "getting number of affected rows")
	}
	if rowsAffected != 1 {
		tx.Rollback()
		return user, errors.Wrap(err, "deleting token")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return user, errors.Wrap(err, "could not commit transaction")
	}

	return user, nil
}

func (u *UserRepo) Update(user User) (User, error) {
	return u.update(context.TODO(), u.db, user)
}

func (u *UserRepo) update(ctx context.Context, tx sqlx.Execer, user User) (User, error) {
	if !user.IsPersisted() {
		return User{}, errors.New("cannot update not persisted record")
	}
	user.Timestamps.JustUpdated(u.clock)

	sql, args, err := u.stmt.Update("users").
		Set("email", user.Email).
		Set("name", user.Name).
		Set("updated_at", user.UpdatedAt).
		Set("password", user.Password).
		Set("last_login", user.LastLogin).
		Set("must_change_password", user.MustChangePassword).
		Where(sq.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return user, errors.Wrap(err, "could not build query")
	}

	result, err := u.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return user, errors.Wrap(err, "could not execute query")
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
	} else if rowsAffected != 1 {
		return user, errors.New("no rows updated")
	}

	return user, nil
}

// Create inserts the given user record into the database. It updates timestamps, and generates a password change token.
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
		Values(user.CreatedAt, user.UpdatedAt, user.Email, user.Password, true).
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

	token, err := u.randomString(32)
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
	cutOff := u.clock().AddDate(0, 0, -1)
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

// ByInviteID finds a user with the given inviteID.
func (u *UserRepo) ByInviteID(ctx context.Context, inviteID string) (User, error) {
	// TODO we need an index on user_password_change_tokens.token
	// TODO we should just add a unique index on it
	sql := "SELECT u.* FROM users u JOIN user_password_change_tokens t on u.id = t.user_id and t.token = $1"
	var user User
	err := u.db.Get(&user, sql, inviteID)
	if err != nil {
		return user, errors.Wrap(err, "could not find user with the given invite id")
	}
	return user, nil
}
