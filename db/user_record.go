package db

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRecord is deprecated
type UserRecord struct {
	Record
	Timestamps

	Email              string     `db:"email"`
	Password           []byte     `db:"password"`
	LastLogin          *time.Time `db:"last_login"`
	MustChangePassword bool       `db:"must_change_password"`
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
