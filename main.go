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

const logFilename = "jornale.txt"

func main() {
	filename := "dictionary.txt"
	// Initialize the log file
	logger, err := middleware.NewLogger(logFilename)
	if err != nil {
		fmt.Println("Error initializing log file:", err)
		return
	}

	d := dictionary.New(filename)

	r := mux.NewRouter()

	r.Use(logger.MiddlewareFunc())
	r.Use(middleware.AuthMiddleware)
	r.HandleFunc("/add", handlers.AddEntryHandler(d, filename)).Methods("POST")
	r.HandleFunc("/get/{word}", handlers.GetDefinitionHandler(d)).Methods("GET")
	r.HandleFunc("/remove/{word}", handlers.RemoveEntryHandler(d)).Methods("DELETE")
	r.HandleFunc("/list", handlers.ListWordsHandler(d)).Methods("GET")

	// Load data from the file
	loadErr := d.LoadFromFile(filename)
	if loadErr != nil {
		fmt.Println("Error loading data:", loadErr)
	}

	http.Handle("/", r)

	fmt.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
