package metrics

import (
	"sync"
	"time"
)

// Metrics - глобальная структура для сбора метрик
type Metrics struct {
	mu sync.RWMutex

	// HTTP metrics
	HTTPRequestsTotal     map[string]int64  // endpoint -> count
	HTTPRequestDuration   map[string][]int64 // endpoint -> durations in ms
	HTTPRequestsInFlight  int64
	HTTPErrorsTotal       map[int]int64 // status_code -> count

	// Database metrics
	DBQueriesTotal        int64
	DBQueryDuration       []int64 // durations in ms
	DBConnectionsOpen     int64
	DBConnectionsIdle     int64
	DBErrorsTotal         int64

	// Cache metrics
	CacheHitsTotal        int64
	CacheMissesTotal      int64
	CacheItemsCount       int64
	CacheEvictionsTotal   int64

	// Application metrics
	StartTime             time.Time
	RequestsTotal         int64
	ErrorsTotal           int64
}

var (
	globalMetrics *Metrics
	once          sync.Once
)

// GetMetrics возвращает singleton instance метрик
func GetMetrics() *Metrics {
	once.Do(func() {
		globalMetrics = &Metrics{
			HTTPRequestsTotal:   make(map[string]int64),
			HTTPRequestDuration: make(map[string][]int64),
			HTTPErrorsTotal:     make(map[int]int64),
			DBQueryDuration:     make([]int64, 0),
			StartTime:           time.Now(),
		}
	})
	return globalMetrics
}

// RecordHTTPRequest записывает метрику HTTP запроса
func (m *Metrics) RecordHTTPRequest(endpoint string, duration time.Duration, statusCode int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.HTTPRequestsTotal[endpoint]++
	m.RequestsTotal++

	durationMs := duration.Milliseconds()
	m.HTTPRequestDuration[endpoint] = append(m.HTTPRequestDuration[endpoint], durationMs)

	// Ограничиваем размер среза до последних 1000 запросов
	if len(m.HTTPRequestDuration[endpoint]) > 1000 {
		m.HTTPRequestDuration[endpoint] = m.HTTPRequestDuration[endpoint][len(m.HTTPRequestDuration[endpoint])-1000:]
	}

	if statusCode >= 400 {
		m.HTTPErrorsTotal[statusCode]++
		m.ErrorsTotal++
	}
}

// IncrementHTTPInFlight увеличивает счетчик текущих запросов
func (m *Metrics) IncrementHTTPInFlight() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.HTTPRequestsInFlight++
}

// DecrementHTTPInFlight уменьшает счетчик текущих запросов
func (m *Metrics) DecrementHTTPInFlight() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.HTTPRequestsInFlight--
}

// RecordDBQuery записывает метрику запроса к БД
func (m *Metrics) RecordDBQuery(duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DBQueriesTotal++
	durationMs := duration.Milliseconds()
	m.DBQueryDuration = append(m.DBQueryDuration, durationMs)

	// Ограничиваем размер среза
	if len(m.DBQueryDuration) > 1000 {
		m.DBQueryDuration = m.DBQueryDuration[len(m.DBQueryDuration)-1000:]
	}

	if err != nil {
		m.DBErrorsTotal++
	}
}

// UpdateDBStats обновляет статистику подключений к БД
func (m *Metrics) UpdateDBStats(open, idle int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DBConnectionsOpen = open
	m.DBConnectionsIdle = idle
}

// RecordCacheHit записывает попадание в кэш
func (m *Metrics) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheHitsTotal++
}

// RecordCacheMiss записывает промах кэша
func (m *Metrics) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheMissesTotal++
}

// UpdateCacheStats обновляет статистику кэша
func (m *Metrics) UpdateCacheStats(items int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheItemsCount = items
}

// RecordCacheEviction записывает вытеснение из кэша
func (m *Metrics) RecordCacheEviction() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheEvictionsTotal++
}

// GetSnapshot возвращает снимок текущих метрик
func (m *Metrics) GetSnapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshot := MetricsSnapshot{
		Uptime:                time.Since(m.StartTime).Seconds(),
		RequestsTotal:         m.RequestsTotal,
		ErrorsTotal:           m.ErrorsTotal,
		HTTPRequestsInFlight:  m.HTTPRequestsInFlight,
		HTTPRequestsTotal:     make(map[string]int64),
		HTTPErrorsTotal:       make(map[int]int64),
		DBQueriesTotal:        m.DBQueriesTotal,
		DBConnectionsOpen:     m.DBConnectionsOpen,
		DBConnectionsIdle:     m.DBConnectionsIdle,
		DBErrorsTotal:         m.DBErrorsTotal,
		CacheHitsTotal:        m.CacheHitsTotal,
		CacheMissesTotal:      m.CacheMissesTotal,
		CacheItemsCount:       m.CacheItemsCount,
		CacheEvictionsTotal:   m.CacheEvictionsTotal,
	}

	// Копируем карты
	for k, v := range m.HTTPRequestsTotal {
		snapshot.HTTPRequestsTotal[k] = v
	}
	for k, v := range m.HTTPErrorsTotal {
		snapshot.HTTPErrorsTotal[k] = v
	}

	// Вычисляем средние значения
	snapshot.HTTPAvgDuration = m.calculateAvgDuration(m.HTTPRequestDuration)
	snapshot.DBAvgDuration = calculateAvg(m.DBQueryDuration)

	// Cache hit rate
	totalCacheRequests := m.CacheHitsTotal + m.CacheMissesTotal
	if totalCacheRequests > 0 {
		snapshot.CacheHitRate = float64(m.CacheHitsTotal) / float64(totalCacheRequests) * 100
	}

	return snapshot
}

// calculateAvgDuration вычисляет среднюю длительность для каждого endpoint
func (m *Metrics) calculateAvgDuration(durations map[string][]int64) map[string]float64 {
	result := make(map[string]float64)
	for endpoint, durs := range durations {
		if len(durs) > 0 {
			result[endpoint] = calculateAvg(durs)
		}
	}
	return result
}

// calculateAvg вычисляет среднее значение
func calculateAvg(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum int64
	for _, v := range values {
		sum += v
	}
	return float64(sum) / float64(len(values))
}

// MetricsSnapshot - снимок метрик в определенный момент времени
type MetricsSnapshot struct {
	// System
	Uptime float64 `json:"uptime_seconds"`

	// HTTP
	RequestsTotal        int64              `json:"requests_total"`
	ErrorsTotal          int64              `json:"errors_total"`
	HTTPRequestsInFlight int64              `json:"http_requests_in_flight"`
	HTTPRequestsTotal    map[string]int64   `json:"http_requests_total"`
	HTTPAvgDuration      map[string]float64 `json:"http_avg_duration_ms"`
	HTTPErrorsTotal      map[int]int64      `json:"http_errors_total"`

	// Database
	DBQueriesTotal    int64   `json:"db_queries_total"`
	DBAvgDuration     float64 `json:"db_avg_duration_ms"`
	DBConnectionsOpen int64   `json:"db_connections_open"`
	DBConnectionsIdle int64   `json:"db_connections_idle"`
	DBErrorsTotal     int64   `json:"db_errors_total"`

	// Cache
	CacheHitsTotal      int64   `json:"cache_hits_total"`
	CacheMissesTotal    int64   `json:"cache_misses_total"`
	CacheHitRate        float64 `json:"cache_hit_rate_percent"`
	CacheItemsCount     int64   `json:"cache_items_count"`
	CacheEvictionsTotal int64   `json:"cache_evictions_total"`
}
