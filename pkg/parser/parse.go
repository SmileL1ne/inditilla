package parser

import (
	"fmt"
	"inditilla/internal/data"
	"inditilla/internal/entity"

	"github.com/dgrijalva/jwt-go/v4"
)

// ParseToken parses given raw token with given signing key and returns
// custom claims extracted from token
func ParseToken(accessToken string, signingKey []byte) (*data.Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &data.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*data.Claims); ok && token.Valid {
		return claims, err
	}

	return nil, entity.ErrInvalidAccessToken
}
