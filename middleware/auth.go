package middleware

import (
	"net/http"
	"strings"

	"github.com/12sub/websockets/auth"
)

func AuthMiddleware(rest http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		_, err := auth.ValidateJsonWebToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		rest(w, r)
	}
}
