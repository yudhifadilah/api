package models

import (
	"github.com/dgrijalva/jwt-go"
)

// jwt struct
type JWTClaims struct {
	jwt.StandardClaims
	IdUser uint `json:"id"`
}
