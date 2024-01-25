// middleware/middleware.go
package middleware

import (
	"net/http"
)

// validToken is the valid authentication token for the middleware.
const validToken = "Bearer hX6FfKOrKHlY7h1rYFxzdtgR7yzzsV0QVmtkH33aDDM="

// AuthMiddleware is a middleware function that checks the Authorization header for a valid token.
// If the token is not valid, it responds with an Unauthorized status.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header.
		token := r.Header.Get("Authorization")

		// Check if the token is valid.
		if token != validToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If the token is valid, proceed to the next handler.
		next.ServeHTTP(w, r)
	})
}
