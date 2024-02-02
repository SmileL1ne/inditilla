package parser

import (
	"fmt"
	"inditilla/internal/entity"
	"inditilla/internal/service/user"

	"github.com/dgrijalva/jwt-go/v4"
)

func ParseToken(accessToken string, signingKey []byte) (*user.Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &user.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*user.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, entity.ErrInvalidAccessToken
}
