package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	Record
	Timestamps

	Email    string `db:"email"`
	Password []byte `db:"password"`
}

func (u *UserRecord) UpdatePassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Printf("Password of user %s updated", u.Email)
	u.Password = hash
	return nil
}

func (u *UserRecord) CheckPassword(compareWith string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(compareWith))
	return err == nil
}
