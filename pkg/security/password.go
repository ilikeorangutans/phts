package security

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// NewPassword creates a new password from the given string.
func NewPassword(passphrase string) (Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "could not crypt password")
	}

	return Password(hash), nil
}

// Password is a thin wrapper around []byte to represent passwords
type Password []byte

// Matches returns true if the given string matches the password.
func (p Password) Matches(compareWith string) bool {
	err := bcrypt.CompareHashAndPassword(p, []byte(compareWith))
	return err == nil
}

// Scan initializes the password from the raw interface{} value, used by the sql driver
func (p *Password) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		b := []byte(v)
		*p = append(*p, b...)
	default:
		return errors.New(fmt.Sprintf("can't parse type %v into password", v))
	}
	return nil
}
