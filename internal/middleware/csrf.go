package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

// CSRFProtection creates CSRF middleware
func CSRFProtection() mux.MiddlewareFunc {
	// Get CSRF auth key from environment
	authKeyStr := os.Getenv("CSRF_AUTH_KEY")
	if authKeyStr == "" {
		log.Fatal("CSRF_AUTH_KEY is not set in environment")
	}

	authKey := []byte(authKeyStr)

	return csrf.Protect(
		authKey,
		csrf.Secure(false),                 // TODO: Set true in production with HTTPS
		csrf.SameSite(csrf.SameSiteLaxMode), // Prevent CSRF attacks
		csrf.Path("/"),
		csrf.MaxAge(24*3600),              // 24 hours to match session
		csrf.ErrorHandler(csrfErrorHandler()), // Custom error handler
		csrf.CookieName("prommsc_csrf"),
		csrf.FieldName("csrf_token"),
		csrf.RequestHeader("X-CSRF-Token"), // For AJAX requests
	)
}

// csrfErrorHandler handles CSRF validation failures
func csrfErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the failed CSRF attempt with IP
		ip := getIP(r)
		log.Printf("CSRF validation failed: IP=%s, Path=%s, Method=%s", ip, r.URL.Path, r.Method)

		http.Error(w, "CSRF token validation failed. Please refresh the page and try again.", http.StatusForbidden)
	})
}

// getIP extracts IP address from request
func getIP(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Use RemoteAddr
	return r.RemoteAddr
}
