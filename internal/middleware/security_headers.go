package middleware

import (
	"net/http"
	"os"
)

// SecurityHeaders добавляет заголовки безопасности ко всем ответам
func SecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// X-Frame-Options: защита от clickjacking
			// DENY - запрещает отображение страницы во фреймах
			w.Header().Set("X-Frame-Options", "DENY")

			// X-Content-Type-Options: предотвращает MIME type sniffing
			// nosniff - браузер должен использовать только объявленный Content-Type
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// X-XSS-Protection: включает встроенную защиту браузера от XSS
			// 1; mode=block - включить защиту и блокировать страницу при обнаружении XSS
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Referrer-Policy: контролирует передачу referrer информации
			// strict-origin-when-cross-origin - передавать только origin при cross-origin запросах
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Content-Security-Policy: защита от XSS и injection атак
			csp := buildCSP()
			w.Header().Set("Content-Security-Policy", csp)

			// Strict-Transport-Security (HSTS): принудительное использование HTTPS
			// Применяется только если SESSION_SECURE=true (production с HTTPS)
			if os.Getenv("SESSION_SECURE") == "true" {
				// max-age=31536000 - 1 год
				// includeSubDomains - применять ко всем поддоменам
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			// Permissions-Policy: контролирует использование браузерных API
			// Отключаем потенциально опасные возможности
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

			next.ServeHTTP(w, r)
		})
	}
}

// buildCSP создает Content-Security-Policy директивы
func buildCSP() string {
	// default-src 'self' - по умолчанию разрешены только ресурсы с того же origin
	// script-src 'self' - скрипты только с того же origin
	// style-src 'self' 'unsafe-inline' - стили с того же origin + inline стили (для tailwind/bootstrap)
	// img-src 'self' data: - изображения с того же origin + data: URLs (для base64)
	// font-src 'self' - шрифты только с того же origin
	// connect-src 'self' - AJAX/WebSocket только к тому же origin
	// frame-ancestors 'none' - запрет на iframe (дублирует X-Frame-Options)
	// base-uri 'self' - ограничение для <base> тега
	// form-action 'self' - формы отправляются только на тот же origin

	return "default-src 'self'; " +
		"script-src 'self'; " +
		"style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data:; " +
		"font-src 'self'; " +
		"connect-src 'self'; " +
		"frame-ancestors 'none'; " +
		"base-uri 'self'; " +
		"form-action 'self'"
}
