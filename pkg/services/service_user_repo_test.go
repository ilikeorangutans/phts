package services

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error mocking connection")
	}
	defer db.Close()

	mock.ExpectBegin()

	repo := NewServiceUsersRepo(sqlx.NewDb(db, "postgres"))

	user, err := repo.NewUser("user@test.local", "secret", false)
	assert.Nil(t, err)
	assert.NotNil(t, user)

}
