package middlewares

import (
	"log"
	"net/http"
	"os"

	"go.mod/security"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFile, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening log file: %v", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)

		fullURL := r.URL.Path
		if r.URL.RawQuery != "" {
			fullURL += "?" + r.URL.RawQuery
		}

		log.Printf("Request: %s %s", r.Method, fullURL)

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	secretHash := os.Getenv("GO_API_SECRET_HASH")
	secretSalt := os.Getenv("GO_API_SECRET_SALT")
	if secretHash == "" || secretSalt == "" {
		log.Fatal("Environment variables GO_API_SECRET_HASH or GO_API_SECRET_SALT are not set.")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if !security.CompareHash(apiKey, secretHash, secretSalt) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
