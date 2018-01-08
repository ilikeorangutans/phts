package model

import "github.com/ilikeorangutans/phts/db"

type User struct {
	db.UserRecord
}

type UserRepository interface {
	FindByEmail(email string) (User, error)
	FindByID(id int64) (User, error)
}

func NewUserRepository(dbx db.DB) UserRepository {
	return &userSQLRepository{
		userDB: db.NewUserDB(dbx),
	}
}

type userSQLRepository struct {
	userDB db.UserDB
}

func (r *userSQLRepository) FindByEmail(email string) (User, error) {
	record, err := r.userDB.FindByEmail(email)

	return User{UserRecord: record}, err
}

func (r *userSQLRepository) FindByID(id int64) (User, error) {
	record, err := r.userDB.FindByID(id)

	return User{UserRecord: record}, err
}
