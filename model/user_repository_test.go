package model

import (
	"testing"

	"github.com/ilikeorangutans/phts/test"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFooBar(t *testing.T) {
	db, mock := test.NewTestDB()

	repository := NewUserRepository(db)
	mock.ExpectExec("INSERT INTO users").WithArgs("test@test.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(123, 1))

	user, err := repository.Create("test@test.com")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	t.Log(mock.ExpectationsWereMet())
	assert.Nil(t, mock.ExpectationsWereMet())
}
