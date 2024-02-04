package data

import (
	"inditilla/pkg/logger"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type TokenModel struct {
	Log *logger.Logger
}

type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func (t *TokenModel) New(email string, deadline time.Duration) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(deadline)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Email: email,
	})

	return token
}
