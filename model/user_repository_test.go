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
	mock.ExpectQuery("INSERT INTO").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint64(13)))

	user, err := repository.Create("test@test.com")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	t.Log(mock.ExpectationsWereMet())
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, int64(13), user.ID)
}
