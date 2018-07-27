package model

import (
	"github.com/ilikeorangutans/phts/db"
)

type UserQueries interface {
	FindByEmail(email string) (*db.UserRecord, error)
	FindByID(id int64) (*db.UserRecord, error)
}

type UserCommands interface {
	UpdatePassword(*db.UserRecord, string) error
	Create(email string) (*db.UserRecord, error)
}

type UserRepository interface {
	UserQueries
	UserCommands
}

func NewUserRepository(dbx db.DB) UserRepository {
	return &userSQLRepository{
		userDB: db.NewUserDB(dbx),
	}
}

type userSQLRepository struct {
	userDB db.UserDB
}

func (r *userSQLRepository) FindByEmail(email string) (*db.UserRecord, error) {
	return r.userDB.FindByEmail(email)
}

func (r *userSQLRepository) FindByID(id int64) (*db.UserRecord, error) {
	return r.userDB.FindByID(id)
}

func (r *userSQLRepository) UpdatePassword(user *db.UserRecord, password string) error {
	if err := user.UpdatePassword(password); err != nil {
		return err
	}

	return r.userDB.Save(user)
}

func (r *userSQLRepository) Create(email string) (user *db.UserRecord, err error) {
	user = &db.UserRecord{
		Email: email,
	}

	randomPassword, err := randomHex(8)
	user.UpdatePassword(randomPassword)
	if err != nil {
		return user, err
	}

	if err := r.userDB.Save(user); err != nil {
		return nil, err
	}
	return user, nil
}
