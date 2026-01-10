# Система мониторинга ProminaYou

## Обзор

ProminaYou включает полноценную систему мониторинга с метриками в формате Prometheus, health checks и детальной телеметрией для отслеживания производительности и состояния приложения.

## Endpoints мониторинга

### 1. Health Check

**Endpoint:** `GET /health`

Проверяет работоспособность всех критичных компонентов приложения.

**Пример запроса:**
```bash
curl http://localhost:8004/health
```

**Пример ответа (healthy):**
```json
{
  "status": "healthy",
  "checks": {
    "database": "ok",
    "memory": "45.23 MB",
    "goroutines": "15"
  }
}
```

**Пример ответа (unhealthy):**
```json
{
  "status": "unhealthy",
  "checks": {
    "database": "failed: connection refused",
    "memory": "512.45 MB",
    "goroutines": "1250"
  }
}
```

**HTTP статус коды:**
- `200 OK` - все системы работают нормально
- `200 OK` (со status: "warning") - есть предупреждения
- `503 Service Unavailable` - система неработоспособна

**Критерии:**
- Database: успешный ping к PostgreSQL
- Memory: потребление памяти (warning если > 500 MB)
- Goroutines: количество горутин (warning если > 1000)

### 2. JSON Metrics

**Endpoint:** `GET /metrics`

Возвращает метрики в формате JSON для простого анализа и интеграции.

**Пример запроса:**
```bash
curl http://localhost:8004/metrics | jq
```

**Пример ответа:**
```json
{
  "uptime_seconds": 3600.5,
  "requests_total": 1523,
  "errors_total": 12,
  "http_requests_in_flight": 3,
  "http_requests_total": {
    "GET /": 450,
    "GET /services": 320,
    "GET /masters": 280,
    "POST /admin/login": 15
  },
  "http_avg_duration_ms": {
    "GET /": 12.5,
    "GET /services": 25.3,
    "GET /masters": 18.7
  },
  "http_errors_total": {
    "404": 8,
    "500": 4
  },
  "db_queries_total": 2340,
  "db_avg_duration_ms": 3.2,
  "db_connections_open": 5,
  "db_connections_idle": 3,
  "db_errors_total": 0,
  "cache_hits_total": 1850,
  "cache_misses_total": 490,
  "cache_hit_rate_percent": 79.05,
  "cache_items_count": 45,
  "cache_evictions_total": 120
}
```

### 3. Prometheus Metrics

**Endpoint:** `GET /metrics/prometheus`

Возвращает метрики в формате Prometheus для интеграции с системами мониторинга.

**Пример запроса:**
```bash
curl http://localhost:8004/metrics/prometheus
```

**Пример ответа:**
```
# HELP prommsc_uptime_seconds Application uptime in seconds
# TYPE prommsc_uptime_seconds gauge
prommsc_uptime_seconds 3600.50

# HELP prommsc_http_requests_total Total HTTP requests
# TYPE prommsc_http_requests_total counter
prommsc_http_requests_total{endpoint="GET /"} 450
prommsc_http_requests_total{endpoint="GET /services"} 320
prommsc_http_requests_total{endpoint="GET /masters"} 280

# HELP prommsc_http_requests_in_flight Current HTTP requests being processed
# TYPE prommsc_http_requests_in_flight gauge
prommsc_http_requests_in_flight 3

# HELP prommsc_http_request_duration_ms Average HTTP request duration in milliseconds
# TYPE prommsc_http_request_duration_ms gauge
prommsc_http_request_duration_ms{endpoint="GET /"} 12.50
prommsc_http_request_duration_ms{endpoint="GET /services"} 25.30

# HELP prommsc_cache_hit_rate Cache hit rate percentage
# TYPE prommsc_cache_hit_rate gauge
prommsc_cache_hit_rate 79.05

# ...и другие метрики
```

## Доступные метрики

### HTTP метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `prommsc_uptime_seconds` | gauge | Время работы приложения |
| `prommsc_http_requests_total` | counter | Общее количество HTTP запросов по endpoint |
| `prommsc_http_requests_in_flight` | gauge | Текущие активные запросы |
| `prommsc_http_request_duration_ms` | gauge | Средняя длительность запроса по endpoint |
| `prommsc_http_errors_total` | counter | Ошибки HTTP по status code |

### Database метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `prommsc_db_queries_total` | counter | Общее количество запросов к БД |
| `prommsc_db_query_duration_ms` | gauge | Средняя длительность запроса к БД |
| `prommsc_db_connections_open` | gauge | Открытые подключения к БД |
| `prommsc_db_connections_idle` | gauge | Простаивающие подключения |
| `prommsc_db_errors_total` | counter | Ошибки при работе с БД |

### Cache метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `prommsc_cache_hits_total` | counter | Попадания в кэш |
| `prommsc_cache_misses_total` | counter | Промахи кэша |
| `prommsc_cache_hit_rate` | gauge | Процент попаданий в кэш |
| `prommsc_cache_items_count` | gauge | Текущее количество элементов в кэше |
| `prommsc_cache_evictions_total` | counter | Вытеснения из кэша |

### Go Runtime метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `prommsc_go_goroutines` | gauge | Количество горутин |
| `prommsc_go_memory_alloc_bytes` | gauge | Выделенная память |
| `prommsc_go_memory_sys_bytes` | gauge | Память от системы |

## Интеграция с Prometheus

### 1. Установка Prometheus

```bash
# Docker
docker run -d -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus

# Или скачать с https://prometheus.io/download/
```

### 2. Конфигурация Prometheus

Создайте файл `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prommsc'
    static_configs:
      - targets: ['localhost:8004']
    metrics_path: '/metrics/prometheus'
    scrape_interval: 10s
```

### 3. Запуск Prometheus

```bash
prometheus --config.file=prometheus.yml
```

### 4. Просмотр метрик

Откройте Prometheus UI: http://localhost:9090

**Примеры запросов (PromQL):**

```promql
# Средняя длительность HTTP запросов
prommsc_http_request_duration_ms

# Rate запросов в секунду
rate(prommsc_http_requests_total[5m])

# Процент ошибок
(sum(rate(prommsc_http_errors_total[5m])) / sum(rate(prommsc_http_requests_total[5m]))) * 100

# Cache hit rate
prommsc_cache_hit_rate

# Использование памяти
prommsc_go_memory_alloc_bytes / 1024 / 1024
```

## Интеграция с Grafana

### 1. Установка Grafana

```bash
docker run -d -p 3000:3000 grafana/grafana
```

### 2. Добавление Prometheus как Data Source

1. Откройте Grafana: http://localhost:3000
2. Login: admin / admin
3. Configuration → Data Sources → Add data source
4. Выберите Prometheus
5. URL: http://localhost:9090
6. Save & Test

### 3. Импорт Dashboard

Создайте dashboard с следующими панелями:

**HTTP Metrics:**
```promql
# Requests per second
rate(prommsc_http_requests_total[5m])

# Average response time
prommsc_http_request_duration_ms

# Error rate
rate(prommsc_http_errors_total[5m])
```

**Database Metrics:**
```promql
# Query rate
rate(prommsc_db_queries_total[5m])

# Connection pool usage
prommsc_db_connections_open

# Query duration
prommsc_db_query_duration_ms
```

**Cache Metrics:**
```promql
# Hit rate
prommsc_cache_hit_rate

# Cache operations
rate(prommsc_cache_hits_total[5m]) + rate(prommsc_cache_misses_total[5m])
```

## Алерты (Alerting)

### Рекомендуемые правила алертов

Создайте файл `alerts.yml`:

```yaml
groups:
  - name: prommsc_alerts
    interval: 30s
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: (sum(rate(prommsc_http_errors_total[5m])) / sum(rate(prommsc_http_requests_total[5m]))) * 100 > 5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }}% (threshold: 5%)"

      # Slow response time
      - alert: SlowResponseTime
        expr: prommsc_http_request_duration_ms > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Slow response time detected"
          description: "Average response time is {{ $value }}ms (threshold: 1000ms)"

      # Low cache hit rate
      - alert: LowCacheHitRate
        expr: prommsc_cache_hit_rate < 50
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Low cache hit rate"
          description: "Cache hit rate is {{ $value }}% (threshold: 50%)"

      # Database connection pool exhaustion
      - alert: DBPoolExhaustion
        expr: prommsc_db_connections_open >= 23
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Database connection pool almost exhausted"
          description: "{{ $value }} of 25 connections in use"

      # Memory usage warning
      - alert: HighMemoryUsage
        expr: prommsc_go_memory_alloc_bytes > 500 * 1024 * 1024
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage is {{ $value | humanize }}B (threshold: 500MB)"

      # Too many goroutines
      - alert: TooManyGoroutines
        expr: prommsc_go_goroutines > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Too many goroutines"
          description: "Goroutine count is {{ $value }} (threshold: 1000)"

      # Service down
      - alert: ServiceDown
        expr: up{job="prommsc"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service is down"
          description: "ProminaYou service is not responding"
```

Добавьте в `prometheus.yml`:
```yaml
rule_files:
  - 'alerts.yml'

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['localhost:9093']
```

## Мониторинг в production

### 1. Автоматизированный мониторинг

```bash
#!/bin/bash
# healthcheck.sh - Скрипт для мониторинга

ENDPOINT="http://localhost:8004/health"
SLACK_WEBHOOK="your-slack-webhook-url"

while true; do
  STATUS=$(curl -s -o /dev/null -w '%{http_code}' $ENDPOINT)

  if [ $STATUS -ne 200 ]; then
    MESSAGE="🚨 ProminaYou health check failed! Status: $STATUS"
    curl -X POST $SLACK_WEBHOOK -H 'Content-Type: application/json' \
      -d "{\"text\": \"$MESSAGE\"}"
  fi

  sleep 60
done
```

### 2. Systemd Service для мониторинга

```ini
[Unit]
Description=ProminaYou Health Check Monitor
After=network.target

[Service]
Type=simple
User=monitor
ExecStart=/usr/local/bin/healthcheck.sh
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 3. Метрики для анализа

**Ключевые показатели (KPIs):**

1. **Availability (Доступность)**
   - Target: > 99.9%
   - Метрика: `up{job="prommsc"}`

2. **Response Time (Время ответа)**
   - Target: < 200ms (p95)
   - Метрика: `prommsc_http_request_duration_ms`

3. **Error Rate (Частота ошибок)**
   - Target: < 1%
   - Метрика: `rate(prommsc_http_errors_total) / rate(prommsc_http_requests_total)`

4. **Throughput (Пропускная способность)**
   - Target: > 100 req/s
   - Метрика: `rate(prommsc_http_requests_total[1m])`

5. **Cache Hit Rate (Эффективность кэша)**
   - Target: > 80%
   - Метрика: `prommsc_cache_hit_rate`

## Troubleshooting

### Проблема: метрики не обновляются

**Решение:**
1. Проверьте, что MetricsMiddleware применен:
```go
r.Use(middleware.MetricsMiddleware())
```

2. Убедитесь, что endpoint доступен:
```bash
curl http://localhost:8004/metrics
```

### Проблема: высокая memory usage

**Диагностика:**
```bash
# Проверка текущего использования памяти
curl http://localhost:8004/metrics | jq '.prommsc_go_memory_alloc_bytes'

# pprof анализ
go tool pprof http://localhost:8004/debug/pprof/heap
```

**Решение:**
- Проверьте размер метрик (ограничено 1000 последних значений)
- Проверьте количество элементов в кэше
- Используйте `runtime.GC()` для принудительной сборки мусора

### Проблема: database connection pool exhaustion

**Диагностика:**
```bash
# Проверка состояния pool
curl http://localhost:8004/metrics | jq '.db_connections_open, .db_connections_idle'
```

**Решение:**
1. Увеличьте MaxOpenConns в config/database.go
2. Проверьте долгие запросы
3. Добавьте timeout для запросов

## Best Practices

1. **Регулярный мониторинг:**
   - Проверяйте /health каждые 30-60 секунд
   - Scrape /metrics/prometheus каждые 10-15 секунд

2. **Алерты:**
   - Настройте критичные алерты (service down, DB unavailable)
   - Настройте warning алерты (slow response, high memory)

3. **Анализ трендов:**
   - Отслеживайте метрики за последние 7-30 дней
   - Выявляйте паттерны использования
   - Планируйте масштабирование заранее

4. **Retention policy:**
   - Храните детальные метрики за последние 7 дней
   - Агрегированные метрики за последние 30 дней
   - Долгосрочные тренды за последний год

## Заключение

Система мониторинга ProminaYou предоставляет полную видимость состояния приложения и позволяет:
- ✅ Отслеживать производительность в реальном времени
- ✅ Быстро обнаруживать и диагностировать проблемы
- ✅ Анализировать тренды и планировать масштабирование
- ✅ Интегрироваться с industry-standard инструментами (Prometheus, Grafana)
