package model

import (
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/pkg/security"
)

type User struct {
	db.Record
	db.Timestamps
	Email               string            `db:"email"`
	Password            security.Password `db:"password"`
	MustChangePassword  bool              `db:"must_change_password"`
	LastLogin           *time.Time        `db:"last_login"`
	PasswordChangeToken string            `db:"-"`
}
