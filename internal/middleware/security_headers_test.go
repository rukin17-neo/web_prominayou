package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestSecurityHeaders проверяет, что все заголовки безопасности устанавливаются корректно
func TestSecurityHeaders(t *testing.T) {
	tests := []struct {
		name           string
		sessionSecure  string
		wantHSTS       bool
		expectedHSTS   string
	}{
		{
			name:          "without HSTS (development)",
			sessionSecure: "false",
			wantHSTS:      false,
		},
		{
			name:          "with HSTS (production)",
			sessionSecure: "true",
			wantHSTS:      true,
			expectedHSTS:  "max-age=31536000; includeSubDomains",
		},
		{
			name:          "without SESSION_SECURE env",
			sessionSecure: "",
			wantHSTS:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменную окружения
			if tt.sessionSecure != "" {
				os.Setenv("SESSION_SECURE", tt.sessionSecure)
			} else {
				os.Unsetenv("SESSION_SECURE")
			}
			defer os.Unsetenv("SESSION_SECURE")

			// Создаем тестовый handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Оборачиваем в SecurityHeaders middleware
			middleware := SecurityHeaders()(handler)

			// Создаем тестовый запрос
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			// Выполняем запрос
			middleware.ServeHTTP(rec, req)

			// Проверяем обязательные заголовки
			headers := map[string]string{
				"X-Frame-Options":        "DENY",
				"X-Content-Type-Options": "nosniff",
				"X-XSS-Protection":       "1; mode=block",
				"Referrer-Policy":        "strict-origin-when-cross-origin",
				"Permissions-Policy":     "geolocation=(), microphone=(), camera=()",
			}

			for headerName, expectedValue := range headers {
				got := rec.Header().Get(headerName)
				if got != expectedValue {
					t.Errorf("%s = %q, want %q", headerName, got, expectedValue)
				}
			}

			// Проверяем Content-Security-Policy
			csp := rec.Header().Get("Content-Security-Policy")
			if csp == "" {
				t.Error("Content-Security-Policy header is missing")
			}

			// Проверяем наличие основных CSP директив
			cspDirectives := []string{
				"default-src 'self'",
				"script-src 'self'",
				"style-src 'self' 'unsafe-inline'",
				"img-src 'self' data:",
				"font-src 'self'",
				"connect-src 'self'",
				"frame-ancestors 'none'",
				"base-uri 'self'",
				"form-action 'self'",
			}

			for _, directive := range cspDirectives {
				if !contains(csp, directive) {
					t.Errorf("CSP missing directive: %s", directive)
				}
			}

			// Проверяем HSTS
			hsts := rec.Header().Get("Strict-Transport-Security")
			if tt.wantHSTS {
				if hsts != tt.expectedHSTS {
					t.Errorf("Strict-Transport-Security = %q, want %q", hsts, tt.expectedHSTS)
				}
			} else {
				if hsts != "" {
					t.Errorf("Strict-Transport-Security should not be set, got %q", hsts)
				}
			}
		})
	}
}

// TestBuildCSP проверяет правильность формирования CSP
func TestBuildCSP(t *testing.T) {
	csp := buildCSP()

	if csp == "" {
		t.Fatal("buildCSP() returned empty string")
	}

	// Проверяем наличие всех необходимых директив
	requiredDirectives := []string{
		"default-src",
		"script-src",
		"style-src",
		"img-src",
		"font-src",
		"connect-src",
		"frame-ancestors",
		"base-uri",
		"form-action",
	}

	for _, directive := range requiredDirectives {
		if !contains(csp, directive) {
			t.Errorf("CSP missing required directive: %s", directive)
		}
	}

	// Проверяем что 'self' используется в большинстве директив
	if !contains(csp, "'self'") {
		t.Error("CSP should contain 'self' keyword")
	}

	// Проверяем что 'none' используется для frame-ancestors
	if !contains(csp, "frame-ancestors 'none'") {
		t.Error("CSP should have 'none' for frame-ancestors")
	}

	// Проверяем что 'unsafe-inline' разрешен только для style-src
	if !contains(csp, "style-src 'self' 'unsafe-inline'") {
		t.Error("CSP should allow 'unsafe-inline' for style-src")
	}

	// Проверяем что data: разрешен для изображений
	if !contains(csp, "img-src 'self' data:") {
		t.Error("CSP should allow data: URIs for images")
	}
}

// TestSecurityHeadersChain проверяет, что middleware корректно передает управление следующему handler
func TestSecurityHeadersChain(t *testing.T) {
	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	middleware := SecurityHeaders()(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if !called {
		t.Error("Next handler was not called")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
}

// contains проверяет содержит ли строка подстроку
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
