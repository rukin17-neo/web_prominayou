package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"prommsc/config"
	"prommsc/internal/metrics"
	"runtime"
)

// MetricsHandler возвращает JSON с метриками
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	m := metrics.GetMetrics()
	snapshot := m.GetSnapshot()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

// PrometheusMetricsHandler возвращает метрики в формате Prometheus
func PrometheusMetricsHandler(w http.ResponseWriter, r *http.Request) {
	m := metrics.GetMetrics()
	snapshot := m.GetSnapshot()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP prommsc_uptime_seconds Application uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE prommsc_uptime_seconds gauge\n")
	fmt.Fprintf(w, "prommsc_uptime_seconds %.2f\n\n", snapshot.Uptime)

	fmt.Fprintf(w, "# HELP prommsc_http_requests_total Total HTTP requests\n")
	fmt.Fprintf(w, "# TYPE prommsc_http_requests_total counter\n")
	for endpoint, count := range snapshot.HTTPRequestsTotal {
		fmt.Fprintf(w, "prommsc_http_requests_total{endpoint=\"%s\"} %d\n", endpoint, count)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP prommsc_http_requests_in_flight Current HTTP requests being processed\n")
	fmt.Fprintf(w, "# TYPE prommsc_http_requests_in_flight gauge\n")
	fmt.Fprintf(w, "prommsc_http_requests_in_flight %d\n\n", snapshot.HTTPRequestsInFlight)

	fmt.Fprintf(w, "# HELP prommsc_http_request_duration_ms Average HTTP request duration in milliseconds\n")
	fmt.Fprintf(w, "# TYPE prommsc_http_request_duration_ms gauge\n")
	for endpoint, duration := range snapshot.HTTPAvgDuration {
		fmt.Fprintf(w, "prommsc_http_request_duration_ms{endpoint=\"%s\"} %.2f\n", endpoint, duration)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP prommsc_http_errors_total Total HTTP errors\n")
	fmt.Fprintf(w, "# TYPE prommsc_http_errors_total counter\n")
	for statusCode, count := range snapshot.HTTPErrorsTotal {
		fmt.Fprintf(w, "prommsc_http_errors_total{status_code=\"%d\"} %d\n", statusCode, count)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "# HELP prommsc_db_queries_total Total database queries\n")
	fmt.Fprintf(w, "# TYPE prommsc_db_queries_total counter\n")
	fmt.Fprintf(w, "prommsc_db_queries_total %d\n\n", snapshot.DBQueriesTotal)

	fmt.Fprintf(w, "# HELP prommsc_db_query_duration_ms Average database query duration in milliseconds\n")
	fmt.Fprintf(w, "# TYPE prommsc_db_query_duration_ms gauge\n")
	fmt.Fprintf(w, "prommsc_db_query_duration_ms %.2f\n\n", snapshot.DBAvgDuration)

	fmt.Fprintf(w, "# HELP prommsc_db_connections_open Open database connections\n")
	fmt.Fprintf(w, "# TYPE prommsc_db_connections_open gauge\n")
	fmt.Fprintf(w, "prommsc_db_connections_open %d\n\n", snapshot.DBConnectionsOpen)

	fmt.Fprintf(w, "# HELP prommsc_db_connections_idle Idle database connections\n")
	fmt.Fprintf(w, "# TYPE prommsc_db_connections_idle gauge\n")
	fmt.Fprintf(w, "prommsc_db_connections_idle %d\n\n", snapshot.DBConnectionsIdle)

	fmt.Fprintf(w, "# HELP prommsc_db_errors_total Total database errors\n")
	fmt.Fprintf(w, "# TYPE prommsc_db_errors_total counter\n")
	fmt.Fprintf(w, "prommsc_db_errors_total %d\n\n", snapshot.DBErrorsTotal)

	fmt.Fprintf(w, "# HELP prommsc_cache_hits_total Total cache hits\n")
	fmt.Fprintf(w, "# TYPE prommsc_cache_hits_total counter\n")
	fmt.Fprintf(w, "prommsc_cache_hits_total %d\n\n", snapshot.CacheHitsTotal)

	fmt.Fprintf(w, "# HELP prommsc_cache_misses_total Total cache misses\n")
	fmt.Fprintf(w, "# TYPE prommsc_cache_misses_total counter\n")
	fmt.Fprintf(w, "prommsc_cache_misses_total %d\n\n", snapshot.CacheMissesTotal)

	fmt.Fprintf(w, "# HELP prommsc_cache_hit_rate Cache hit rate percentage\n")
	fmt.Fprintf(w, "# TYPE prommsc_cache_hit_rate gauge\n")
	fmt.Fprintf(w, "prommsc_cache_hit_rate %.2f\n\n", snapshot.CacheHitRate)

	fmt.Fprintf(w, "# HELP prommsc_cache_items_count Current cache items count\n")
	fmt.Fprintf(w, "# TYPE prommsc_cache_items_count gauge\n")
	fmt.Fprintf(w, "prommsc_cache_items_count %d\n\n", snapshot.CacheItemsCount)

	// Go runtime metrics
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Fprintf(w, "# HELP prommsc_go_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE prommsc_go_goroutines gauge\n")
	fmt.Fprintf(w, "prommsc_go_goroutines %d\n\n", runtime.NumGoroutine())

	fmt.Fprintf(w, "# HELP prommsc_go_memory_alloc_bytes Bytes allocated and still in use\n")
	fmt.Fprintf(w, "# TYPE prommsc_go_memory_alloc_bytes gauge\n")
	fmt.Fprintf(w, "prommsc_go_memory_alloc_bytes %d\n\n", mem.Alloc)

	fmt.Fprintf(w, "# HELP prommsc_go_memory_sys_bytes Total bytes from system\n")
	fmt.Fprintf(w, "# TYPE prommsc_go_memory_sys_bytes gauge\n")
	fmt.Fprintf(w, "prommsc_go_memory_sys_bytes %d\n\n", mem.Sys)
}

// HealthCheckHandler проверяет здоровье приложения
func HealthCheckHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := HealthStatus{
			Status: "healthy",
			Checks: make(map[string]string),
		}

		// Проверка базы данных
		if err := db.Ping(); err != nil {
			health.Status = "unhealthy"
			health.Checks["database"] = "failed: " + err.Error()
		} else {
			health.Checks["database"] = "ok"

			// Получение статистики БД
			stats := config.GetDatabaseStats(db)
			openConns := stats["open_connections"].(int)
			m := metrics.GetMetrics()
			m.UpdateDBStats(int64(openConns), int64(stats["idle"].(int)))
		}

		// Проверка памяти
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		health.Checks["memory"] = fmt.Sprintf("%.2f MB", float64(mem.Alloc)/1024/1024)

		// Проверка goroutines
		goroutines := runtime.NumGoroutine()
		health.Checks["goroutines"] = fmt.Sprintf("%d", goroutines)
		if goroutines > 1000 {
			health.Status = "warning"
		}

		// Отправка ответа
		w.Header().Set("Content-Type", "application/json")
		if health.Status == "unhealthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else if health.Status == "warning" {
			w.WriteHeader(http.StatusOK) // 200, но со статусом warning
		}

		json.NewEncoder(w).Encode(health)
	}
}

// HealthStatus представляет статус здоровья приложения
type HealthStatus struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}
