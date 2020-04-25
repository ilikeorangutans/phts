package services

import (
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/security"
)

type ServiceUser struct {
	db.Record
	db.Timestamps
	Email     string            `db:"email"`
	Password  security.Password `db:"password"`
	LastLogin *time.Time        `db:"last_login"`
	// SystemCreated indicates that this record has been automatically created and changes to it will likely be overwritten when the app restarts.
	// This is only the case for the system admin user.
	SystemCreated bool `db:"system_created"`
}

// CheckPassword returns true if the given string equals the password of the user.
func (s ServiceUser) CheckPassword(compareWith string) bool {
	return s.Password.Matches(compareWith)
}
