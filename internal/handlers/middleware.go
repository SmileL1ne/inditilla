package handlers

import (
	"context"
	"inditilla/pkg/parser"
	"net/http"
	"os"
	"strings"
)

const contextKey = "user"

func (r *routes) jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			r.l.Error("Unauthorized")
			// SEND ERROR RESPNONSE HERE
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			r.l.Error("Unauthorized")
			// SEND ERROR RESPNONSE HERE
			return
		}

		if headerParts[0] != "Bearer" {
			r.l.Error("Unauthorized")
			// SEND ERROR RESPNONSE HERE
			return
		}

		claims, err := parser.ParseToken(headerParts[1], []byte(os.Getenv("SIGNING_KEY")))

		if claims == nil {
			if err != nil {
				r.l.Error(err.Error())
				// SEND ERROR RESPNONSE HERE
				return
			}
			r.l.Error("Unauthorized")
			// SEND ERROR RESPNONSE HERE
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), contextKey, claims))

		next.ServeHTTP(w, req)
	})
}
