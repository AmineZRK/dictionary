// logging.go
package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Logger struct encapsulates the log-related functionality.
type Logger struct {
	LogFile *os.File
	Logger  *log.Logger
}

// LoggingMiddlewareFunc is a function signature for the LoggingMiddleware.
type LoggingMiddlewareFunc func(http.Handler) http.Handler

// NewLogger initializes a new Logger instance.
func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := log.New(file, "", log.LstdFlags)

	return &Logger{LogFile: file, Logger: logger}, nil
}

// LoggingMiddleware adds logging functionality to HTTP requests.
func (l *Logger) LoggingMiddleware() LoggingMiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			l.Logger.Printf("[%s] %s %s %s\n", time.Since(start), r.Method, r.RequestURI, r.RemoteAddr)
		})
	}
}

// SetLogFile creates or opens a log file for logging.
func (l *Logger) SetLogFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	l.LogFile = file
	l.Logger.SetOutput(l.LogFile)

	return nil
}

// MiddlewareFunc converts LoggingMiddlewareFunc to mux.MiddlewareFunc.
func (l *Logger) MiddlewareFunc() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return l.LoggingMiddleware()(next)
	}
}
