package model

import "github.com/ilikeorangutans/phts/db"

type User struct {
	//Timestamps
	ID     int64
	Handle string
	Email  string
}

type UserRepository interface {
	FindByID(id uint) (User, error)
	FindByHandle(handle string) (User, error)
}

type userSQLRepository struct {
	db db.DB
}
