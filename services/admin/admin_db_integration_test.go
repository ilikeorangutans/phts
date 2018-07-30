package admin

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
	dbtest "github.com/ilikeorangutans/phts/test/integration/db"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewAdminRecord(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		adminDB := NewAdminDB(dbx)
		user, _ := dbtest.CreateUser(t, dbx)

		record := &AdminRecord{
			UserRecord: *user,
		}
		adminDB.Save(record)

		assert.Equal(t, record.UserID, record.ID)
		assert.True(t, record.UserID > 0)
	})
}

func TestFindByEmail(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		_, user, adminDB := CreateAdmin(t, dbx)

		admin, err := adminDB.FindByEmail(user.Email)
		assert.Nil(t, err)
		assert.NotNil(t, admin)
		assert.Equal(t, user.Email, admin.Email)
	})
}

func TestFindByEmailAndPasswordReturnsErrWithWrongPassword(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		_, user, _ := CreateAdmin(t, dbx)
		service := NewAdminService(dbx)

		admin, err := service.FindByEmailAndPassword(user.Email, "wrong password")

		assert.NotNil(t, err)
		assert.Nil(t, admin)
	})
}

func CreateAdmin(t *testing.T, dbx db.DB) (*AdminRecord, *db.UserRecord, AdminDB) {
	user, _ := dbtest.CreateUser(t, dbx)
	adminDB := NewAdminDB(dbx)

	record := &AdminRecord{
		UserRecord: *user,
		UserID:     user.ID,
	}

	assert.Nil(t, adminDB.Save(record))
	return record, user, adminDB
}
