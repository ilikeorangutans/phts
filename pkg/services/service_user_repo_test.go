package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func withServiceUserRepo(t *testing.T, f func(mock sqlmock.Sqlmock, repo *ServiceUsersRepo, now time.Time)) {
	now := time.Now()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error mocking connection")
	}
	defer db.Close()

	repo := NewServiceUsersRepo(sqlx.NewDb(db, "postgres"))
	repo.clock = func() time.Time { return now }
	repo.newPassword = func(input string) (security.Password, error) { return []byte(fmt.Sprintf("%s encrypted", input)), nil }

	f(mock, repo, now)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Logf("unfulfilled expectations: %+v", err)
	}
}

func TestNewUser(t *testing.T) {
	withServiceUserRepo(t, func(mock sqlmock.Sqlmock, repo *ServiceUsersRepo, now time.Time) {
		mock.
			ExpectQuery("insert into service_users").
			WithArgs("user@test.local", []byte("secret encrypted"), now, now, true).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(13))

		user, err := repo.NewUser("user@test.local", "secret", true)

		assert.Nil(t, err)
		assert.Equal(t, int64(13), user.ID)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
		assert.True(t, user.SystemCreated)
	})
}

func TestUpdatePassword(t *testing.T) {
	withServiceUserRepo(t, func(mock sqlmock.Sqlmock, repo *ServiceUsersRepo, now time.Time) {
		mock.
			ExpectExec("UPDATE service_users .* WHERE id = ").
			WithArgs([]byte("updated password encrypted"), now, nil, false, 13).
			WillReturnResult(sqlmock.NewResult(0, 1))

		user := ServiceUser{
			Record: db.Record{
				ID: 13,
			},
			Email:    "",
			Password: []byte("secret encrypted"),
		}

		user, err := repo.UpdatePassword(user, "updated password")

		assert.Nil(t, err)
		assert.Equal(t, now, user.UpdatedAt)
	})
}

func TestJustLoggedIn(t *testing.T) {
	withServiceUserRepo(t, func(mock sqlmock.Sqlmock, repo *ServiceUsersRepo, now time.Time) {
		yesterday := now.AddDate(0, 0, -1)
		mock.
			ExpectExec("UPDATE service_users .* WHERE id = ").
			WithArgs([]byte(""), now, now, false, 13).
			WillReturnResult(sqlmock.NewResult(0, 1))

		user := ServiceUser{
			Record: db.Record{
				ID: 13,
			},
			LastLogin: &yesterday,
			Password:  []byte(""),
		}

		user, err := repo.JustLoggedIn(user)

		assert.Nil(t, err)
		assert.Equal(t, now, user.UpdatedAt)
	})
}

func TestFindByEmail(t *testing.T) {
	withServiceUserRepo(t, func(mock sqlmock.Sqlmock, repo *ServiceUsersRepo, now time.Time) {
		mock.
			ExpectQuery("^SELECT .* FROM service_users WHERE email .* LIMIT 1$").
			WithArgs("foo@test.local").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "test.foo"))

		_, err := repo.FindByEmail("foo@test.local")

		assert.Nil(t, err)
	})
}
