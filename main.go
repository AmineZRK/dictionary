// main.go
package main

import (
	"estiam/dictionary"
	"estiam/handlers"
	"estiam/middleware"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// logFilename is the name of the log file.
const logFilename = "jornale.txt"

func main() {
	// Initialize the logger for logging middleware.
	logger, err := middleware.NewLogger(logFilename)
	if err != nil {
		fmt.Println("Error initializing log file:", err)
		return
	}

	// Initialize the dictionary.
	d, err := dictionary.NewDictionary("mongodb://localhost:27017", "dictionary", "dictionary")
	if err != nil {
		fmt.Println("Error initializing dictionary:", err)
		return
	}

	// Create a new Gorilla Mux router.
	r := mux.NewRouter()

	// Use the logger middleware for logging requests.
	r.Use(logger.MiddlewareFunc())

	// Use the AuthMiddleware for authentication.
	r.Use(middleware.AuthMiddleware)

	// Define routes and corresponding handlers.
	r.HandleFunc("/add", handlers.AddEntryHandler(d)).Methods("POST")
	r.HandleFunc("/get/{word}", handlers.GetDefinitionHandler(d)).Methods("GET")
	r.HandleFunc("/remove/{word}", handlers.RemoveEntryHandler(d)).Methods("DELETE")
	r.HandleFunc("/list", handlers.ListWordsHandler(d)).Methods("GET")

	// Set up the HTTP server with the Gorilla Mux router.
	http.Handle("/", r)

	// Start the server on port 8080.
	fmt.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
