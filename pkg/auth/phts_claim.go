package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type PhtsClaim struct {
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"email"`
	jwt.StandardClaims
}
