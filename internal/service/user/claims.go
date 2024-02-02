package user

import "github.com/dgrijalva/jwt-go/v4"

type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}
