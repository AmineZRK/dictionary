// main.go
package main

import (
	"encoding/json"
	"estiam/dictionary"
	"estiam/middleware"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

const logFilename = "jornale.txt"

func main() {
	filename := "dictionary.txt"
	// Initialize the log file
	errs := middleware.SetLogFile(logFilename)
	if errs != nil {
		fmt.Println("Error initializing log file:", errs)
		return
	}

	d := dictionary.New(filename)

	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/add", addEntryHandler(d, filename)).Methods("POST")
	r.HandleFunc("/get/{word}", getDefinitionHandler(d)).Methods("GET")
	r.HandleFunc("/remove/{word}", removeEntryHandler(d)).Methods("DELETE")
	r.HandleFunc("/list", listWordsHandler(d)).Methods("GET")

	// Load data from the file
	err := d.LoadFromFile(filename)
	if err != nil {
		fmt.Println("Error loading data:", err)
	}

	http.Handle("/", r)

	fmt.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}

func addEntryHandler(d *dictionary.Dictionary, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entry dictionary.EntryOperation
		err := json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
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

func getDefinitionHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		word := params["word"]

		entry, err := d.Get(word)
		if err != nil {
			http.Error(w, "Word not found", http.StatusAccepted)
			return
		}

		response := map[string]string{"word": word, "definition": entry.Definition}
		jsonResponse(w, response)
	}
}

func removeEntryHandler(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		word := params["word"]

		message, err := d.Remove(word)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error removing word: %v", err), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, map[string]string{"message": message})
	}
}

func listWordsHandler(d *dictionary.Dictionary) http.HandlerFunc {
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
