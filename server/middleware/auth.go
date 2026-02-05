package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const IsAdminKey contextKey = "isAdmin"

func APIKeyAuthMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			expected := "Bearer " + apiKey
			if auth != expected {
				http.Error(w, "You are unauthorized to access the requested resource.", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func OptionalAdminMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			expected := "Bearer " + apiKey
			isAdmin := auth == expected
			ctx := context.WithValue(r.Context(), IsAdminKey, isAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
