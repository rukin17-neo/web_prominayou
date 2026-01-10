# ProminaYou - Massage Salon Management System

[![Go Version](https://img.shields.io/badge/Go-1.24.3-00ADD8?logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-60%2B%20passing-success)](TESTING.md)
[![Coverage](https://img.shields.io/badge/Coverage-100%25%20(core)-success)](TESTING.md)

> Полнофункциональное веб-приложение для управления массажным салоном с административной панелью, системой мониторинга и оптимизацией производительности.

## 🌟 Ключевые возможности

### Для клиентов
- 📋 Просмотр услуг и мастеров
- 💬 Чтение отзывов
- 📞 Контактная информация
- 📱 Адаптивный дизайн

### Для администраторов
- 👥 Управление мастерами (CRUD)
- 💆 Управление услугами (CRUD)
- 🔐 Безопасная аутентификация
- 📊 Система мониторинга
- 🚀 Оптимизация производительности

## 🔥 Новые функции (v1.0)

### 📊 Система мониторинга
- **Health checks** - `/health` endpoint для проверки состояния приложения
- **JSON metrics** - `/metrics` для интеграции с custom dashboards
- **Prometheus metrics** - `/metrics/prometheus` для Prometheus/Grafana
- **Real-time телеметрия** - HTTP, Database, Cache, Runtime метрики

### ⚡ Оптимизация производительности
- **Database indexes** - автоматическое создание индексов при старте
- **In-memory cache** - SimpleCache с TTL для query results
- **Template caching** - предварительная загрузка и кэширование шаблонов
- **Connection pooling** - оптимизированный пул подключений к БД
- **pprof profiling** - endpoints для CPU/memory профилирования

### 🔒 Безопасность
- **Security headers** - CSP, X-Frame-Options, HSTS, etc.
- **Rate limiting** - защита от брутфорса и DoS атак
- **Input validation** - RFC-compliant email/password validation
- **Error handling** - безопасная обработка ошибок без утечки информации
- **Timing attack protection** - константное время для password checks

### ✅ Тестирование
- **60+ unit tests** - comprehensive test coverage
- **100% coverage** - критичные модули (validation, pagination, security)
- **Benchmarks** - производительность валидации и пагинации

## 📋 Требования

- Go 1.24.3 или выше
- PostgreSQL 16
- Git

## 🚀 Быстрый старт

### 1. Клонирование репозитория

```bash
git clone https://github.com/rukin17-neo/web_prominayou.git
cd web_prominayou
```

### 2. Установка зависимостей

```bash
go mod download
```

### 3. Настройка базы данных

```sql
-- Создание базы данных
CREATE DATABASE salon;

-- Создание пользователя
CREATE USER salon_admin WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE salon TO salon_admin;
```

### 4. Конфигурация окружения

Скопируйте `.env.example` в `.env` и настройте:

```bash
cp .env.example .env
```

Отредактируйте `.env`:
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=salon_admin
DB_PASSWORD=your_password
DB_NAME=salon
DB_SSLMODE=disable

# Session
SESSION_SECURE=false  # true для production
SESSION_SECRET=generate-random-32-byte-string

# SMTP (для восстановления пароля)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Application
APP_URL=http://localhost:8004
```

**Генерация секретов:**
```bash
# SESSION_SECRET
openssl rand -base64 32

# CSRF_AUTH_KEY (если включена CSRF защита)
openssl rand -base64 32
```

### 5. Запуск приложения

```bash
go run cmd/web/main.go
```

Приложение будет доступно по адресу: **http://localhost:8004**

## 🗂️ Структура проекта

```
prommsc/
├── cmd/web/
│   └── main.go                    # Entry point
├── config/
│   ├── database.go                # Database configuration
│   ├── email.go                   # SMTP configuration
│   └── optimization.go            # DB optimization
├── internal/
│   ├── auth/                      # Authentication
│   ├── cache/                     # In-memory cache
│   ├── handlers/
│   │   ├── admin/                 # Admin panel handlers
│   │   ├── client/                # Public handlers
│   │   └── shared/                # Shared utilities
│   ├── metrics/                   # Prometheus metrics
│   └── middleware/                # HTTP middleware
├── models/                        # Data models & repositories
├── templates/                     # HTML templates
├── static/                        # CSS, images
├── MONITORING.md                  # Monitoring guide
├── PERFORMANCE.md                 # Performance guide
├── TESTING.md                     # Testing guide
└── prometheus.yml                 # Prometheus config
```

## 🌐 API Endpoints

### 🔓 Публичные endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/` | Главная страница |
| GET | `/services` | Список услуг с пагинацией |
| GET | `/masters` | Список мастеров с пагинацией |
| GET | `/masters/photo/{id}` | Фотография мастера |
| GET | `/contacts` | Контакты |
| GET | `/reviews` | Отзывы |

### 🔐 Административная панель

| Метод | Endpoint | Описание | Rate Limit |
|-------|----------|----------|------------|
| GET | `/admin/login` | Страница входа | - |
| POST | `/admin/login` | Вход в систему | 5/15мин |
| POST | `/admin/logout` | Выход | - |
| GET/POST | `/admin/forgot-password` | Восстановление пароля | 3/15мин |
| GET/POST | `/admin/reset-password` | Сброс пароля | 5/15мин |
| GET/POST | `/admin/masters` | CRUD мастеров | 30/мин |
| GET/POST | `/admin/services` | CRUD услуг | 30/мин |
| GET/POST | `/admin/users` | CRUD пользователей | 20/мин |

### 📊 Мониторинг и метрики

| Метод | Endpoint | Описание | Формат |
|-------|----------|----------|--------|
| GET | `/health` | Health check | JSON |
| GET | `/metrics` | Метрики приложения | JSON |
| GET | `/metrics/prometheus` | Prometheus метрики | Text |
| GET | `/debug/pprof/` | CPU/Memory profiling | HTML/Binary |

**Пример health check:**
```bash
curl http://localhost:8004/health
```
```json
{
  "status": "healthy",
  "checks": {
    "database": "ok",
    "memory": "2.42 MB",
    "goroutines": "7"
  }
}
```

**Пример метрик:**
```bash
curl http://localhost:8004/metrics | jq
```
```json
{
  "uptime_seconds": 60.97,
  "requests_total": 150,
  "http_requests_total": {
    "GET /": 45,
    "GET /services": 32
  },
  "db_connections_open": 5,
  "cache_hit_rate_percent": 85.5
}
```

## 📊 Интеграция с Prometheus

### Конфигурация Prometheus

```yaml
scrape_configs:
  - job_name: 'prommsc'
    static_configs:
      - targets: ['localhost:8004']
    metrics_path: '/metrics/prometheus'
    scrape_interval: 10s
```

### Доступные метрики

- **HTTP:** requests, duration, errors, in-flight
- **Database:** queries, connections, errors
- **Cache:** hits, misses, hit rate, items
- **Runtime:** goroutines, memory

Подробнее: [MONITORING.md](MONITORING.md)

## ⚡ Производительность

### Оптимизации

- ✅ Database indexes (masters, services, users, sessions)
- ✅ In-memory cache с TTL
- ✅ Template caching
- ✅ Connection pooling (25 max connections)
- ✅ Query optimization с ANALYZE

### Benchmarks

```
BenchmarkEmailValidation       2,952,188 ops    408 ns/op    0 allocs
BenchmarkPasswordValidation   31,425,751 ops     33 ns/op    0 allocs
BenchmarkPagination         1,000,000,000 ops    0.3 ns/op    0 allocs
BenchmarkCache_Set             10,000,000 ops    120 ns/op    1 allocs
BenchmarkCache_Get             20,000,000 ops     60 ns/op    0 allocs
```

Подробнее: [PERFORMANCE.md](PERFORMANCE.md)

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты
go test ./...

# С покрытием
go test -cover ./...

# Verbose режим
go test -v ./...

# Benchmarks
go test -bench=. -benchmem ./...
```

### Coverage

- ✅ **validation.go** - 100%
- ✅ **pagination.go** - 100%
- ✅ **security_headers.go** - 100%
- ✅ **errors.go** - 100%

Подробнее: [TESTING.md](TESTING.md)

## 🔒 Безопасность

### Security Headers

```
✅ X-Frame-Options: DENY
✅ X-Content-Type-Options: nosniff
✅ X-XSS-Protection: 1; mode=block
✅ Referrer-Policy: strict-origin-when-cross-origin
✅ Content-Security-Policy: (9 директив)
✅ Permissions-Policy: geolocation=(), microphone=(), camera=()
✅ Strict-Transport-Security: (только HTTPS)
```

### Rate Limiting

| Категория | Лимит | Период |
|-----------|-------|--------|
| Login | 5 | 15 мин |
| Password Reset | 3/5 | 15 мин |
| User Management | 20 | 1 мин |
| CRUD Operations | 30 | 1 мин |

### Input Validation

- ✅ Email: RFC 5322 compliant
- ✅ Password: 8-72 chars, uppercase + lowercase + digits
- ✅ Pagination: max limit 100 (DoS protection)

## 📚 Документация

- **[MONITORING.md](MONITORING.md)** - Полное руководство по мониторингу (521 строка)
- **[PERFORMANCE.md](PERFORMANCE.md)** - Оптимизация производительности (387 строк)
- **[TESTING.md](TESTING.md)** - Тестирование и QA (559 строк)
- **[.env.example](.env.example)** - Пример конфигурации

## 🚢 Deployment

### Production Checklist

- [ ] Сгенерировать новые SESSION_SECRET и CSRF_AUTH_KEY
- [ ] Установить SESSION_SECURE=true
- [ ] Настроить HTTPS (для HSTS)
- [ ] Изменить DB_PASSWORD на сложный
- [ ] Настроить реальный SMTP
- [ ] Обновить APP_URL на production домен
- [ ] Настроить Prometheus scraping
- [ ] Создать Grafana dashboards
- [ ] Настроить алерты в AlertManager

### Docker (опционально)

```dockerfile
FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o prommsc cmd/web/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/prommsc .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
EXPOSE 8004
CMD ["./prommsc"]
```

## 🤝 Разработка

### Commit Convention

```
git commit -m "feat: add user profile page"
git commit -m "fix: resolve login redirect issue"
git commit -m "docs: update API documentation"
git commit -m "test: add pagination tests"
```

### Pull Requests

1. Fork репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'feat: add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Создайте Pull Request

## 📊 Статистика проекта

- **Язык:** Go 1.24.3
- **Строк кода:** ~4,000+
- **Документация:** ~1,500 строк
- **Тесты:** 60+
- **Coverage:** 100% (критичные модули)
- **Go файлов:** 48
- **Commits:** 15+
- **Pull Requests:** 2 merged

## 📄 Лицензия

MIT License - см. [LICENSE](LICENSE)

## 👨‍💻 Автор

**rukin17-neo**
- GitHub: [@rukin17-neo](https://github.com/rukin17-neo)
- Email: rukin.a0889@gmail.com

## 🙏 Благодарности

- Разработка при участии **Claude Sonnet 4.5** ([Claude Code](https://claude.com/claude-code))
- PostgreSQL community
- Go community
- Prometheus & Grafana teams

---

⭐ Если проект оказался полезным, поставьте звезду на GitHub!
