// middleware_test.go
package middleware_test

import (
	"estiam/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareValidToken(t *testing.T) {
	// 1. Create a new http.ResponseWriter and http.Request.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/list", nil)
	assert.NoError(t, err)

	// 2. Set a valid token in the request header.
	r.Header.Set("Authorization", "Bearer hX6FfKOrKHlY7h1rYFxzdtgR7yzzsV0QVmtkH33aDDM=")

	// 3. Create a handler using middleware.AuthMiddleware.
	handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 4. Serve an HTTP request using the handler.
	handler.ServeHTTP(w, r)

	// 5. Verify that the request is successful (status code 200).
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	// 1. Create a new http.ResponseWriter and http.Request.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// 2. Set an invalid token in the request header.
	r.Header.Set("Authorization", "Bearer invalid_token")

	// 3. Create a handler using middleware.AuthMiddleware.
	handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 4. Serve an HTTP request using the handler.
	handler.ServeHTTP(w, r)

	// 5. Verify that the response has a status code of 401 (Unauthorized).
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestValidateDataValid(t *testing.T) {
	// 1. Call middleware.ValidateData with valid word and definition.
	err := middleware.ValidateData("valid_word", "valid_definition")

	// 2. Verify that the function returns no error.
	assert.NoError(t, err)
}

func TestValidateDataInvalid(t *testing.T) {
	// 1. Call middleware.ValidateData with invalid word and definition.
	err := middleware.ValidateData("sh", "def") // Invalid data

	// 2. Verify that the function returns an error.
	assert.Error(t, err)
}

func TestHandleError(t *testing.T) {
	// 1. Create a new http.ResponseWriter.
	w := httptest.NewRecorder()

	// 2. Call middleware.HandleError with a sample error message and status code.
	middleware.HandleError(w, "Sample error", http.StatusInternalServerError)

	// 3. Verify that the response has the expected status code and body.
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Sample error")
}
