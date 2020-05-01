package model

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ilikeorangutans/phts/pkg/security"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func withServiceUserRepo(t *testing.T, f func(mock sqlmock.Sqlmock, repo *UserRepo, now time.Time)) {
	now := time.Now()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error mocking connection")
	}
	defer db.Close()

	repo := NewUserRepo(sqlx.NewDb(db, "postgres"))
	repo.clock = func() time.Time { return now }
	repo.newPassword = func(input string) (security.Password, error) { return []byte(fmt.Sprintf("%s encrypted", input)), nil }
	repo.randomString = func(n int) (string, error) { return strings.Repeat("x", n), nil }

	f(mock, repo, now)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Logf("unfulfilled expectations: %+v", err)
	}
}

func TestNewUser(t *testing.T) {
	withServiceUserRepo(t, func(mock sqlmock.Sqlmock, repo *UserRepo, now time.Time) {

		mock.ExpectExec("^DELETE FROM user_password_change_tokens WHERE created_at .*$").
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO users").
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(13))
		mock.ExpectExec("^INSERT INTO user_password_change_tokens").
			WithArgs(13, now, strings.Repeat("x", 32), true).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		user, err := repo.NewUser("test@test.local")
		assert.Nil(t, err)
		assert.NotNil(t, user)
	})
}

func TestNewUserRollsBack(t *testing.T) {
	withServiceUserRepo(t, func(mock sqlmock.Sqlmock, repo *UserRepo, now time.Time) {

		mock.ExpectExec("^DELETE FROM user_password_change_tokens WHERE created_at .*$").
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO users").
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(13))
		mock.ExpectExec("^INSERT INTO user_password_change_tokens").
			WithArgs(13, now, strings.Repeat("x", 32), true).
			WillReturnError(errors.New("stuff broke"))
		mock.ExpectRollback()

		_, err := repo.NewUser("test@test.local")
		assert.NotNil(t, err)
	})
}
