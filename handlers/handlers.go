// handlers.go
package handlers

import (
	"encoding/json"
	"estiam/dictionary"
	"estiam/middleware"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// AddEntryHandler is a handler for adding entries to the dictionary.
func AddEntryHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the incoming JSON request into an EntryOperation.
		var entry dictionary.EntryOperation
		err := json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the incoming data.
		err = middleware.ValidateData(entry.Word, entry.Definition)
		if err != nil {
			middleware.HandleError(w, fmt.Sprintf("Error validating data: %v", err), http.StatusBadRequest)
			return
		}

		fmt.Printf("Received JSON: %+v\n", entry)

		// Trim leading and trailing whitespaces from word and definition.
		word := strings.TrimSpace(entry.Word)
		definition := strings.TrimSpace(entry.Definition)

		// Add the word to the dictionary.
		message, err := d.Add(word, definition)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding word: %v", err), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, map[string]string{"message": message})
	}
}

// GetDefinitionHandler retrieves the definition of a word from the dictionary.
func GetDefinitionHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the word parameter from the request.
		params := mux.Vars(r)
		word := params["word"]

		// Get the definition of the word from the dictionary.
		entry, err := d.Get(word)
		if err != nil {
			middleware.HandleError(w, fmt.Sprintf("Error getting word: %v", err), http.StatusNotFound)
			return
		}

		// Prepare and send the response.
		response := map[string]string{"word": word, "definition": entry.String()}
		jsonResponse(w, response)
	}
}

// RemoveEntryHandler removes a word and its definition from the dictionary.
func RemoveEntryHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the word parameter from the request.
		params := mux.Vars(r)
		word := params["word"]

		// Remove the word from the dictionary.
		message, err := d.Remove(word)
		if err != nil {
			middleware.HandleError(w, fmt.Sprintf("Error removing word: %v", err), http.StatusInternalServerError)
			return
		}

		// Prepare and send the response.
		jsonResponse(w, map[string]string{"message": message})
	}
}

// ListWordsHandler retrieves a list of all words in the dictionary.
func ListWordsHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the list of words from the dictionary.
		words, _ := d.List()

		// Prepare and send the response.
		response := map[string][]string{"words": words}
		jsonResponse(w, response)
	}
}

// jsonResponse sets the Content-Type header to JSON and encodes the given data as JSON.
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
