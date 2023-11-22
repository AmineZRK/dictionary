package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

var logFile *os.File

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s %s %s\n", time.Since(start), r.Method, r.RequestURI, r.RemoteAddr)
	})
}

func SetLogFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logFile = file
	log.SetOutput(logFile)

	return nil
}
