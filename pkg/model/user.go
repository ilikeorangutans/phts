package model

import (
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/security"
)

// User is a admin user in phts.
type User struct {
	db.Record
	db.Timestamps
	Email               string            `db:"email"`
	Password            security.Password `db:"password"`
	MustChangePassword  bool              `db:"must_change_password"`
	LastLogin           *time.Time        `db:"last_login"`
	PasswordChangeToken string            `db:"-"`
	Name                string            `db:"name"`
}

// UserFromOldRecord is to help transition from the old style records to the new, simpler ones
func UserFromOldRecord(user *db.UserRecord) User {
	return User{
		Record:     user.Record,
		Timestamps: user.Timestamps,
		Email:      user.Email,
	}
}
