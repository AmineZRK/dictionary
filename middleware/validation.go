// middleware/validation.go
package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func ValidateData(word, definition string) error {
	// Validate data length
	if len(word) < 3 || len(definition) < 5 {
		return fmt.Errorf("invalid data: Word and definition must be at least 3 and 5 characters long, respectively")
	}

	// Validate no special characters in word and definition
	if containsSpecialCharacters(word) || containsSpecialCharacters(definition) {
		return fmt.Errorf("invalid data: Word and definition must not contain special characters")
	}

	return nil
}

// checks if a string contains special characters
func containsSpecialCharacters(s string) bool {
	specialCharacters := "~!@#$%^&*()-+={}[]|;:'\",.<>?/"
	for _, char := range s {
		if strings.ContainsRune(specialCharacters, char) {
			return true
		}
	}
	return false
}

// handles errors by logging and sending an appropriate HTTP response
func HandleError(w http.ResponseWriter, message string, statusCode int) {
	fmt.Println("Error:", message)
	http.Error(w, message, statusCode)
}
