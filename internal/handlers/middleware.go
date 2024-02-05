package handlers

import (
	"context"
	"errors"
	"fmt"
	"inditilla/pkg/parser"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const contextKey = "user"

// jwtAuth is a middleware that authenticates user by given jwt token.
// It returns 401 Status Unauthorized if no token given or it is invalid
//
// (for just viewing resource, it may be considered to set user as 'anonymous'
// if no token is provided)
func (r *routes) jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Vary", "Authorization")

		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		/* Additionally, may let user in as anonymous user here */

		// If token is present check and validate it
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		if isValidToken := r.validateToken(headerParts[1]); !isValidToken {
			r.invalidAuthToken(w, req, "Authentication")
			return
		}

		// Parse token with signing key from environment file
		claims, err := parser.ParseToken(headerParts[1], []byte(os.Getenv("SIGNING_KEY")))
		if err != nil {
			r.l.Error("jwtAuth: %v", err)
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		// Check if such user exists
		exists, err := r.s.User.Exists(req.Context(), claims.Email)
		if !exists {
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					r.invalidAuthToken(w, req, "Authentication")
					return
				}
				r.serverError(w, req, err, "Authentcation")
				return
			}
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		// Check if token is not expired
		if time.Since(claims.ExpiresAt.Time) >= 0 {
			r.invalidAuthToken(w, req, "Authentcation")
			return
		}

		// Put claims to request's context by custom context key
		req = req.WithContext(context.WithValue(req.Context(), contextKey, claims))
		next.ServeHTTP(w, req)
	})
}

func (r *routes) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				r.serverError(w, req, fmt.Errorf("%s", err), "recover panic")
			}
		}()

		next.ServeHTTP(w, req)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self';")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, req)
	})
}
