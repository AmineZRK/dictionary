package middleware

import (
	"net/http"
)

const validToken = "Bearer hX6FfKOrKHlY7h1rYFxzdtgR7yzzsV0QVmtkH33aDDM="

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != validToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
