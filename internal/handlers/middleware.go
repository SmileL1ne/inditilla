package handlers

import (
	"context"
	"errors"
	"inditilla/internal/entity"
	"inditilla/pkg/parser"
	"net/http"
	"os"
	"strings"
	"time"
)

const contextKey = "user"

func (r *routes) jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Vary", "Authorization")

		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		// Additionally, may let user in as anonymous user here

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		claims, err := parser.ParseToken(headerParts[1], []byte(os.Getenv("SIGNING_KEY")))
		if claims == nil {
			if errors.Is(err, entity.ErrInvalidAccessToken) {
				r.invalidAuthToken(w, req, "Authentcation")
				return
			}
			r.serverError(w, req, err, "Authentcation")
			return
		}

		exists, err := r.s.User.Exists(req.Context(), claims.Email)
		if !exists {
			if err != nil {
				r.serverError(w, req, err, "Authentcation")
				return
			}
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		if time.Since(claims.ExpiresAt.Time) >= 0 {
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), contextKey, claims))
		next.ServeHTTP(w, req)
	})
}
