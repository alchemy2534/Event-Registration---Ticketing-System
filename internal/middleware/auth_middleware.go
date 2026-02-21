package middleware

import (
	"log"
	"net/http"
)

// LoggingMiddleware logs all incoming HTTP requests to help track traffic.
// As time permits, this could be expanded to an AuthMiddleware by injecting
// context and validating headers before passing it to `next.ServeHTTP`.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
