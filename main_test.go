// main_test.go
package main_test

import (
	"bytes"
	"encoding/json"
	"estiam/dictionary"
	"estiam/handlers"
	"estiam/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// 1. Create a new dictionary and logger.
	d := dictionary.New("test_dictionary.txt")
	logger, err := middleware.NewLogger("test_log.txt")
	if err != nil {
		t.Fatal("Error creating logger:", err)
	}

	// 2. Create a new http.ResponseWriter and http.Request for the add endpoint.
	w := httptest.NewRecorder()
	entry := dictionary.EntryOperation{
		Word:       "test_word",
		Definition: "test_definition",
	}
	entryJSON, err := json.Marshal(entry)
	if err != nil {
		t.Fatal("Error encoding JSON:", err)
	}

	req, err := http.NewRequest("POST", "/add", bytes.NewBuffer(entryJSON))
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// 3. Create a handler using AddEntryHandler.
	handler := handlers.AddEntryHandler(d, "test_dictionary.txt")

	// 4. Create a router and apply the logger middleware directly.
	r := mux.NewRouter()
	r.Use(logger.MiddlewareFunc())
	r.HandleFunc("/add", handler).Methods("POST")

	// 5. Serve an HTTP request using the router.
	r.ServeHTTP(w, req)

	// 6. Verify the response status code.
	assert.Equal(t, http.StatusOK, w.Code, "Status code should be OK")

	// 7. Verify the response body or any other expectations based on your implementation.
	// Example: Check if the response contains a JSON message.
	expectedMessage := `{"message":"Word 'test_word' Added successfully"}`
	assert.Equal(t, strings.TrimSpace(expectedMessage), strings.TrimSpace(w.Body.String()), "Response body should match expected message")
}

func TestGetDefinitionHandler(t *testing.T) {
	// 1. Create a new dictionary and logger.
	d := dictionary.New("test_dictionary.txt")
	logger, err := middleware.NewLogger("test_log.txt")
	if err != nil {
		t.Fatal("Error creating logger:", err)
	}

	// 2. Create a new http.ResponseWriter and http.Request for the get endpoint.
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/get/test_word", nil)
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// 3. Create a handler using GetDefinitionHandler.
	handler := handlers.GetDefinitionHandler(d)

	// 4. Create a router and apply the logger middleware directly.
	r := mux.NewRouter()
	r.Use(logger.MiddlewareFunc())
	r.HandleFunc("/get/{word}", handler).Methods("GET")

	// 5. Serve an HTTP request using the router.
	r.ServeHTTP(w, req)

	// 6. Verify the response status code.
	assert.Equal(t, http.StatusAccepted, w.Code, "Status code should be Accepted")

	// 7. Verify the response body or any other expectations based on your implementation.
	// Example: Check if the response contains a JSON message.
	expectedResponse := "Error getting word: word not found: test_word\n"
	assert.Equal(t, expectedResponse, w.Body.String(), "Response body should match expected response")
}

func TestRemoveEntryHandler(t *testing.T) {
	// 1. Create a new dictionary and logger.
	d := dictionary.New("test_dictionary.txt")
	logger, err := middleware.NewLogger("test_log.txt")
	if err != nil {
		t.Fatal("Error creating logger:", err)
	}

	// Add a word to the dictionary for testing removal
	d.Add("test_word", "test_definition")

	// 2. Create a new http.ResponseWriter and http.Request for the remove endpoint.
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/remove/test_word", nil)
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// 3. Create a handler using RemoveEntryHandler.
	handler := handlers.RemoveEntryHandler(d)

	// 4. Create a router and apply the logger middleware directly.
	r := mux.NewRouter()
	r.Use(logger.MiddlewareFunc())
	r.HandleFunc("/remove/{word}", handler).Methods("DELETE")

	// 5. Serve an HTTP request using the router.
	r.ServeHTTP(w, req)

	// 6. Verify the response status code.
	assert.Equal(t, http.StatusOK, w.Code, "Status code should be OK")

	// 7. Verify the response body or any other expectations based on your implementation.
	// Example: Check if the response contains a JSON message.
	expectedResponse := `{"message":"Word 'test_word' removed successfully"}`
	assert.Contains(t, w.Body.String(), expectedResponse, "Response body should contain expected response")
}

func TestListWordsHandler(t *testing.T) {
	// 1. Create a new dictionary and logger.
	d := dictionary.New("test_dictionary.txt")
	logger, err := middleware.NewLogger("test_log.txt")
	if err != nil {
		t.Fatal("Error creating logger:", err)
	}

	// Add words to the dictionary for testing listing
	d.Add("test_word_1", "definition_1")
	d.Add("test_word_2", "definition_2")

	// 2. Create a new http.ResponseWriter and http.Request for the list endpoint.
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/list", nil)
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// 3. Create a handler using ListWordsHandler.
	handler := handlers.ListWordsHandler(d)

	// 4. Create a router and apply the logger middleware directly.
	r := mux.NewRouter()
	r.Use(logger.MiddlewareFunc())
	r.HandleFunc("/list", handler).Methods("GET")

	// 5. Serve an HTTP request using the router.
	r.ServeHTTP(w, req)

	// 6. Verify the response status code.
	assert.Equal(t, http.StatusOK, w.Code, "Status code should be OK")

	// 7. Verify the response body or any other expectations based on your implementation.
	// Example: Check if the response contains a JSON message.
	actualResponse := w.Body.String()

	// Check if all expected words are present in the actual response
	for _, word := range []string{"test_word_1", "test_word_2"} {
		assert.Contains(t, actualResponse, word, "Response body should contain expected word")
	}
}
