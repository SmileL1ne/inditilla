package handlers

import (
	"context"
	"fmt"
	"inditilla/pkg/parser"
	"net/http"
	"os"
	"strings"
	"time"
)

const contextKey = "user"

func (r *routes) jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			r.l.Error("Unauthorized")
			// SEND ERROR RESPONSE HERE
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			r.l.Error("Unauthorized")
			// SEND ERROR RESPONSE HERE
			return
		}

		if headerParts[0] != "Bearer" {
			r.l.Error("Unauthorized")
			// SEND ERROR RESPONSE HERE
			return
		}

		claims, err := parser.ParseToken(headerParts[1], []byte(os.Getenv("SIGNING_KEY")))
		if claims == nil {
			if err != nil {
				r.l.Error(err.Error())
				// SEND SERVER ERROR RESPONSE HERE
				return
			}
			r.l.Error("Unauthorized")
			// SEND ERROR RESPONSE HERE
			return
		}

		exists, err := r.s.User.Exists(req.Context(), claims.Email)
		if !exists {
			if err != nil {
				r.l.Error(err.Error())
				// SEND SERVER ERROR RESPONSE HERE
				return
			}
			r.l.Error("Unauthorized")
			// SEND ERROR RESPONSE HERE
			return
		}

		if time.Since(claims.ExpiresAt.Time) >= 0 {
			fmt.Println("nah")
			r.l.Error("Unauthorized")
			// SEND ERROR RESPONSE HERE
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), contextKey, claims))
		next.ServeHTTP(w, req)
	})
}
