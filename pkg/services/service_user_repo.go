package services

import (
	godb "database/sql"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/ilikeorangutans/phts/pkg/security"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ServiceUsersPaginator are the default paginator options for service users.
var ServiceUsersPaginator = database.OffsetPaginatorOpts{
	MinLimit:           1,
	DefaultLimit:       10,
	MaxLimit:           100,
	ValidOrderColumns:  []string{"id", "created_at", "updated_at", "email"},
	DefaultOrderColumn: "id",
	DefaultOrder:       "asc",
}

func NewServiceUsersRepo(db *sqlx.DB) *ServiceUsersRepo {
	return &ServiceUsersRepo{
		db:          db,
		clock:       time.Now,
		newPassword: security.NewPassword,
		stmt:        sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type ServiceUsersRepo struct {
	db          *sqlx.DB
	clock       func() time.Time
	newPassword func(string) (security.Password, error)
	stmt        sq.StatementBuilderType
}

func (s *ServiceUsersRepo) List(paginator database.OffsetPaginator) ([]ServiceUser, database.OffsetPaginator, error) {
	var users []ServiceUser
	var count uint64
	err := s.stmt.RunWith(s.db).Select("count(*)").From("service_users").QueryRow().Scan(&count)
	if err != nil {
		return users, paginator, errors.Wrap(err, "could not list service users")
	}
	paginator = paginator.WithCount(count)

	query := s.stmt.Select("*").From("service_users")
	query = paginator.Paginate(query)

	sql, args, err := query.ToSql()
	if err != nil {
		return users, paginator, errors.Wrap(err, "could not list service users")
	}

	err = s.db.Select(&users, sql, args...)
	if err != nil {
		return users, paginator, errors.Wrap(err, "could not list service users")
	}

	return users, paginator, nil
}

// JustLoggedIn updates the given user's LastLogin field to now
func (s *ServiceUsersRepo) JustLoggedIn(user ServiceUser) (ServiceUser, error) {
	now := s.clock()
	user.LastLogin = &now

	return s.Update(user)
}

// FindByEmail finds a user by email
func (s *ServiceUsersRepo) FindByEmail(email string) (ServiceUser, error) {
	var result ServiceUser

	sql, args, err := s.stmt.Select("*").
		From("service_users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return result, errors.Wrap(err, "could not build query")
	}

	if err := s.db.Get(&result, sql, args...); err == godb.ErrNoRows {
		return result, err
	} else if err != nil {
		return result, errors.Wrap(err, "could not select record")
	}

	return result, nil
}

// NewUser creates a new user, persists it in the database, and returns the created record.
func (s *ServiceUsersRepo) NewUser(email, password string, system bool) (ServiceUser, error) {
	p, err := s.newPassword(password)
	if err != nil {
		return ServiceUser{}, errors.Wrap(err, "could not create admin user")
	}
	user := ServiceUser{SystemCreated: system, Email: email, Password: p}
	return s.Create(user)
}

// UpdatePassword updates the given users' password
func (s *ServiceUsersRepo) UpdatePassword(user ServiceUser, password string) (ServiceUser, error) {
	p, err := s.newPassword(password)
	if err != nil {
		return user, errors.Wrap(err, "could not crypt password")
	}
	user.Password = p

	return s.Update(user)
}

func (s *ServiceUsersRepo) Update(user ServiceUser) (ServiceUser, error) {
	user.UpdatedAt = s.clock()

	sql, args, err := s.stmt.
		Update("service_users").
		Set("password", user.Password).
		Set("updated_at", user.UpdatedAt).
		Set("last_login", user.LastLogin).
		Set("must_change_password", user.MustChangePassword).
		Where(sq.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return user, errors.Wrap(err, "could not generate sql")
	}

	result, err := s.db.Exec(sql, args...)
	if err != nil {
		return user, errors.Wrap(err, "could not update user record")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return user, errors.Wrap(err, "could not check number of rows affected")
	}
	if rows != 1 {
		return user, errors.New("update didn't update one row")
	}

	return user, nil
}

// Create inserts a new user into the database.
func (s *ServiceUsersRepo) Create(user ServiceUser) (ServiceUser, error) {
	user.Timestamps = db.JustCreated(s.clock)

	row := s.db.QueryRow("insert into service_users (email, password, created_at, updated_at, system_created) values ($1, $2, $3, $4, $5) returning id", user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.SystemCreated)
	err := row.Scan(&user.ID)
	if err != nil {
		return user, errors.Wrap(err, "could not create user record")
	}
	return user, nil
}
