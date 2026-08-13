package middleware

import (
	"net/http"
	"strings"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/api/response"
)

func Auth(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			if parts[1] != token {
				response.Error(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
