// handlers/handlers.go
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

const logFilename = "jornale.txt"

// AddEntryHandler is a handler for adding entries to the dictionary.
func AddEntryHandler(d *dictionary.Dictionary, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entry dictionary.EntryOperation
		err := json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = middleware.ValidateData(entry.Word, entry.Definition)
		if err != nil {
			middleware.HandleError(w, fmt.Sprintf("Error validating data: %v", err), http.StatusBadRequest)
			return
		}

		fmt.Printf("Received JSON: %+v\n", entry)

		word := strings.TrimSpace(entry.Word)
		definition := strings.TrimSpace(entry.Definition)

		message, err := d.Add(word, definition)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding word: %v", err), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, map[string]string{"message": message})
		//w.WriteHeader(http.StatusCreated)
	}
}

// GetDefinitionHandler is a handler for getting the definition of a word.
func GetDefinitionHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		word := params["word"]

		entry, err := d.Get(word)
		if err != nil {
			middleware.HandleError(w, fmt.Sprintf("Error getting word: %v", err), http.StatusAccepted)
			return
		}

		response := map[string]string{"word": word, "definition": entry.Definition}
		jsonResponse(w, response)
	}
}

// RemoveEntryHandler is a handler for removing entries from the dictionary.
func RemoveEntryHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		word := params["word"]

		message, err := d.Remove(word)
		if err != nil {
			middleware.HandleError(w, fmt.Sprintf("Error removing word: %v", err), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, map[string]string{"message": message})
	}
}

// ListWordsHandler is a handler for listing all words in the dictionary.
func ListWordsHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		words, _ := d.List()

		response := map[string][]string{"words": words}
		jsonResponse(w, response)
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
