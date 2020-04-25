package services

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/session"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticationHandler(sessions session.Storage, usersRepo *ServiceUsersRepo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO clear any existing sessions
		// TODO disregard session ids coming from the client

		defer r.Body.Close()

		if err := r.ParseForm(); err != nil {
			http.Error(w, "could not parse form", http.StatusBadRequest)
			return
		}

		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		user, err := usersRepo.FindByEmail(email)
		if err != nil {
			// TODO add error message to request
			LoginHandler(w, r)
			return
		}

		if !user.CheckPassword(password) {
			// TODO add error message to request
			LoginHandler(w, r)
			return
		}

		sessionID, err := GenerateRandomString(32)
		if err != nil {
			http.Error(w, "could not generate random string", http.StatusInternalServerError)
			return
		}
		sessions.Add(sessionID, nil)

		http.Redirect(w, r, "/services/internal/", http.StatusFound)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login page")
}

type ServiceUser struct {
	db.Record
	db.Timestamps
	Email     string    `db:"email"`
	Password  []byte    `db:"password"`
	LastLogin time.Time `db:"last_login"`
}

// CheckPassword returns true if the given string equals the password of the user.
func (s ServiceUser) CheckPassword(compareWith string) bool {
	err := bcrypt.CompareHashAndPassword(s.Password, []byte(compareWith))
	return err == nil
}

func NewServiceUsersRepo(db db.DB) *ServiceUsersRepo {
	return &ServiceUsersRepo{
		db: db,
	}
}

type ServiceUsersRepo struct {
	db db.DB
}

func (s *ServiceUsersRepo) FindByEmail(email string) (ServiceUser, error) {
	var result ServiceUser

	if err := s.db.Get(&result, "select * from service_users where email = $1", email); err != nil {
		return result, errors.Wrap(err, "could not select record")
	}

	return result, nil
}

func (s *ServiceUsersRepo) UpdatePassword(user ServiceUser, password string) (ServiceUser, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, errors.Wrap(err, "could not crypt password")
	}
	user.Password = hash

	return s.Update(user)
}

func (s *ServiceUsersRepo) Update(user ServiceUser) (ServiceUser, error) {
	user.UpdatedAt = time.Now()

	result, err := s.db.Exec("update service_users set password = $1, updated_at = $1 where email = $3", user.Password, user.UpdatedAt, user.Email)
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

	result, err := s.db.Exec("insert into service_users (email, password, created_at, updated_at) values ($1, $2, $3, $4) returning id", user.Password, user.UpdatedAt, user.Email)
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

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
