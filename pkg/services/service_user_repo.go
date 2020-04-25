package services

import (
	godb "database/sql"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/pkg/errors"
)

func NewServiceUsersRepo(db db.DB) *ServiceUsersRepo {
	return &ServiceUsersRepo{
		db:    db,
		clock: time.Now,
	}
}

type ServiceUsersRepo struct {
	db    db.DB
	clock func() time.Time
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

	result, err := s.db.Exec("update service_users set password = $1, updated_at = $2 where email = $3", user.Password, user.UpdatedAt, user.Email)
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
	now := time.Now()
	user.UpdatedAt = now
	user.CreatedAt = now

	row := s.db.QueryRow("insert into service_users (email, password, created_at, updated_at, system_created) values ($1, $2, $3, $4, $5) returning id", user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.SystemCreated)
	err := row.Scan(&user.ID)
	if err != nil {
		return user, errors.Wrap(err, "could not create user record")
	}
	return user, nil
}
