package middleware

import (
	"net/http"
	"prommsc/internal/metrics"
	"time"

	"github.com/gorilla/mux"
)

// responseWriter обертка для захвата status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // default 200
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware собирает метрики HTTP запросов
func MetricsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			m := metrics.GetMetrics()

			// Увеличиваем счетчик активных запросов
			m.IncrementHTTPInFlight()
			defer m.DecrementHTTPInFlight()

			// Оборачиваем ResponseWriter для захвата status code
			rw := newResponseWriter(w)

			// Выполняем запрос
			next.ServeHTTP(rw, r)

			// Записываем метрики
			duration := time.Since(start)
			endpoint := r.Method + " " + r.URL.Path
			m.RecordHTTPRequest(endpoint, duration, rw.statusCode)
		})
	}
}
