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

	csrfMiddleware := csrf.Protect(
		authKey,
		csrf.Secure(false),                  // Отключаем Secure для localhost
		csrf.SameSite(csrf.SameSiteStrictMode), // Strict mode
		csrf.Path("/"),
		csrf.MaxAge(24*3600),                // 24 hours to match session
		csrf.ErrorHandler(csrfErrorHandler()), // Custom error handler
		csrf.CookieName("csrf_token_v2"),    // Новое имя cookie
		csrf.FieldName("csrf_token"),
		csrf.RequestHeader("X-CSRF-Token"), // For AJAX requests
	)

	log.Printf("CSRF middleware initialized: Secure=false, SameSite=Strict, MaxAge=24h")
	return csrfMiddleware
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
// Защита от IP spoofing: использует rightmost IP из X-Forwarded-For
func getIP(r *http.Request) string {
	// Check X-Forwarded-For header (берем последний IP - ближайший к серверу)
	// Это защищает от spoofing, так как клиент может подделать только начало списка
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		// Берем последний IP (rightmost) вместо первого
		return strings.TrimSpace(parts[len(parts)-1])
	}

	// Check X-Real-IP header (может быть установлен только доверенным proxy)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Use RemoteAddr (прямое соединение)
	// Удаляем порт если есть
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}
