package services

import (
	godb "database/sql"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/database"
	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
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
		db:    db,
		clock: time.Now,
	}
}

type ServiceUsersRepo struct {
	db    *sqlx.DB
	clock func() time.Time
}

func (s *ServiceUsersRepo) List(paginator database.OffsetPaginator) ([]ServiceUser, database.OffsetPaginator, error) {
	var users []ServiceUser
	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	var count uint64
	err := stmt.RunWith(s.db).Select("count(*)").From("service_users").QueryRow().Scan(&count)
	if err != nil {
		return users, paginator, errors.Wrap(err, "could not list service users")
	}
	paginator = paginator.WithCount(count)

	query := stmt.Select("*").From("service_users")
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

func (s *ServiceUsersRepo) JustLoggedIn(user ServiceUser) (ServiceUser, error) {
	now := time.Now()
	user.LastLogin = &now

	return s.Update(user)
}

func (s *ServiceUsersRepo) FindByEmail(email string) (ServiceUser, error) {
	var result ServiceUser

	if err := s.db.Get(&result, "select * from service_users where email = $1", email); err == godb.ErrNoRows {
		return result, err
	} else if err != nil {
		return result, errors.Wrap(err, "could not select record")
	}

	return result, nil
}

// NewUser creates a new user, persists it in the database, and returns the created record.
func (s *ServiceUsersRepo) NewUser(email, password string, system bool) (ServiceUser, error) {
	p, err := security.NewPassword(password)
	if err != nil {
		return ServiceUser{}, errors.Wrap(err, "could not create admin user")
	}
	user := ServiceUser{SystemCreated: system, Email: email, Password: p}
	return s.Create(user)
}

func (s *ServiceUsersRepo) UpdatePassword(user ServiceUser, password string) (ServiceUser, error) {
	p, err := security.NewPassword(password)
	if err != nil {
		return user, errors.Wrap(err, "could not crypt password")
	}
	user.Password = p

	return user, nil
}

func (s *ServiceUsersRepo) Update(user ServiceUser) (ServiceUser, error) {
	user.UpdatedAt = time.Now()

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := stmt.
		Update("service_users").
		Set("password", user.Password).
		Set("updated_at", user.UpdatedAt).
		Set("last_login", user.LastLogin).
		Set("must_change_password", user.MustChangePassword).
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
