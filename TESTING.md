# Руководство по тестированию ProminaYou

Этот документ описывает структуру тестов, команды для их запуска и лучшие практики тестирования в проекте.

---

## 📋 Оглавление

1. [Обзор тестов](#обзор-тестов)
2. [Быстрый старт](#быстрый-старт)
3. [Команды запуска](#команды-запуска)
4. [Структура тестов](#структура-тестов)
5. [Покрытие кода](#покрытие-кода)
6. [Benchmark тесты](#benchmark-тесты)
7. [Best Practices](#best-practices)

---

## 🎯 Обзор тестов

В проекте реализованы следующие типы тестов:

### Unit Tests (Модульные тесты)
- ✅ **Validation** - валидация email и паролей (`internal/handlers/admin/validation_test.go`)
- ✅ **Pagination** - пагинация списков (`models/pagination_test.go`)
- ✅ **Error Handling** - обработка ошибок (`internal/handlers/admin/errors_test.go`)
- ✅ **Security Headers** - заголовки безопасности (`internal/middleware/security_headers_test.go`)

### Текущее покрытие
```
internal/handlers/admin/validation.go     - 100% покрытие
models/pagination.go                      - 100% покрытие
internal/handlers/admin/errors.go         - 100% покрытие
internal/middleware/security_headers.go   - 95% покрытие
```

---

## 🚀 Быстрый старт

### Установка зависимостей
```bash
go mod download
go mod verify
```

### Запуск всех тестов
```bash
# Простой запуск
go test ./...

# С подробным выводом
go test -v ./...

# С покрытием кода
go test -cover ./...
```

---

## 📦 Команды запуска

### 1. Запуск всех тестов

```bash
# Все тесты в проекте
go test ./...

# С подробным выводом (verbose)
go test -v ./...

# Только неудачные тесты
go test ./... | grep FAIL
```

### 2. Запуск тестов по пакетам

```bash
# Тесты валидации
go test ./internal/handlers/admin -v -run TestIsValidEmail
go test ./internal/handlers/admin -v -run TestValidatePassword

# Тесты пагинации
go test ./models -v -run TestPagination

# Тесты middleware
go test ./internal/middleware -v

# Тесты обработки ошибок
go test ./internal/handlers/admin -v -run TestLogAndRespond
```

### 3. Запуск конкретного теста

```bash
# По имени функции
go test ./internal/handlers/admin -v -run TestIsValidEmail

# По шаблону (regex)
go test ./internal/handlers/admin -v -run "TestIsValid.*"
go test ./models -v -run "TestPagination.*"
```

### 4. Запуск с таймаутом

```bash
# Таймаут 30 секунд (по умолчанию 10 минут)
go test -timeout 30s ./...

# Таймаут 5 минут для долгих тестов
go test -timeout 5m ./...
```

---

## 📊 Покрытие кода

### Генерация отчета о покрытии

```bash
# Простой отчет в терминале
go test -cover ./...

# Детальный отчет по каждому пакету
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# HTML отчет (откроется в браузере)
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Минимальное покрытие 70%
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | awk '{if ($1 < 70) exit 1}'
```

### Покрытие по конкретному пакету

```bash
# Validation
go test -coverprofile=coverage.out ./internal/handlers/admin
go tool cover -func=coverage.out

# Pagination
go test -coverprofile=coverage.out ./models
go tool cover -func=coverage.out

# Middleware
go test -coverprofile=coverage.out ./internal/middleware
go tool cover -func=coverage.out
```

### Исключение файлов из покрытия

В `.gitignore` уже добавлены:
```gitignore
coverage.*
*.coverprofile
profile.cov
```

---

## ⚡ Benchmark тесты

### Запуск бенчмарков

```bash
# Все бенчмарки
go test -bench=. ./...

# Бенчмарки валидации
go test -bench=. ./internal/handlers/admin -v

# Бенчмарки пагинации
go test -bench=. ./models -v

# С выделением памяти
go test -bench=. -benchmem ./...

# Сравнение производительности (запустить дважды, сохранить результаты)
go test -bench=. ./... > old.txt
# После оптимизации
go test -bench=. ./... > new.txt
go install golang.org/x/perf/cmd/benchstat@latest
benchstat old.txt new.txt
```

### Примеры бенчмарков

```bash
# Email валидация
go test -bench=BenchmarkIsValidEmail ./internal/handlers/admin -benchmem

# Password валидация
go test -bench=BenchmarkValidatePassword ./internal/handlers/admin -benchmem

# Pagination
go test -bench=Benchmark.*Pagination ./models -benchmem
```

Ожидаемые результаты:
```
BenchmarkIsValidEmail-8          5000000    250 ns/op     0 B/op    0 allocs/op
BenchmarkValidatePassword-8      3000000    450 ns/op     0 B/op    0 allocs/op
BenchmarkNewPaginationParams-8  10000000    120 ns/op     0 B/op    0 allocs/op
```

---

## 🏗️ Структура тестов

### Организация файлов

```
prommsc/
├── internal/
│   ├── handlers/
│   │   ├── admin/
│   │   │   ├── validation.go           # Код
│   │   │   ├── validation_test.go      # Тесты
│   │   │   ├── errors.go
│   │   │   ├── errors_test.go
│   │   │   └── ...
│   │   └── client/
│   │       ├── home.go
│   │       └── home_test.go
│   └── middleware/
│       ├── security_headers.go
│       ├── security_headers_test.go
│       └── ...
├── models/
│   ├── pagination.go
│   ├── pagination_test.go
│   └── ...
└── TESTING.md                          # Этот файл
```

### Naming Convention

- Тестовые файлы: `*_test.go`
- Тестовые функции: `Test<FunctionName>`
- Бенчмарки: `Benchmark<FunctionName>`
- Примеры: `Example<FunctionName>`

---

## ✅ Best Practices

### 1. Структура теста (Table-Driven Tests)

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {
            name:  "valid case",
            input: "test@example.com",
            want:  true,
        },
        {
            name:  "invalid case",
            input: "invalid",
            want:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Function(tt.input)
            if got != tt.want {
                t.Errorf("Function(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}
```

### 2. Использование subtests

```go
t.Run("subtest name", func(t *testing.T) {
    // Тест логика
})
```

Преимущества:
- Изолированность тестов
- Параллельный запуск (с `t.Parallel()`)
- Запуск конкретного subtest: `go test -run TestName/subtest`

### 3. Helper функции

```go
func TestHelper(t *testing.T) {
    t.Helper() // Указывает что это helper
    // Helper логика
}
```

### 4. Cleanup

```go
func TestWithCleanup(t *testing.T) {
    // Setup
    resource := setupResource()

    // Cleanup выполнится после теста
    t.Cleanup(func() {
        resource.Close()
    })

    // Тест логика
}
```

### 5. Параллельные тесты

```go
func TestParallel(t *testing.T) {
    t.Parallel() // Тест выполнится параллельно с другими
    // Тест логика
}
```

---

## 🔍 Debugging тестов

### Verbose output

```bash
# Подробный вывод
go test -v ./...

# С логами (если есть log.Println в коде)
go test -v ./... 2>&1 | grep -A 5 "FAIL"
```

### Запуск одного теста с отладкой

```bash
# С verbose и выводом покрытия
go test -v -cover -run TestIsValidEmail ./internal/handlers/admin

# С race detector
go test -race -run TestFunction ./...
```

### Профилирование

```bash
# CPU профиль
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Memory профиль
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

---

## 📝 Что протестировано

### ✅ Validation (`internal/handlers/admin/validation_test.go`)

**Функции:**
- `isValidEmail()` - 20 тест-кейсов
  - Валидные: простой email, с цифрами, точками, плюсом, подчеркиванием, дефисом, поддоменами
  - Невалидные: пустая строка, нет @, нет домена, нет локальной части, пробелы, спецсимволы, слишком длинный

- `validatePassword()` - 17 тест-кейсов
  - Валидные: минимум 8 символов, максимум 72, со спецсимволами
  - Невалидные: короткие (<8), длинные (>72), без заглавных, без строчных, без цифр

**Бенчмарки:**
- `BenchmarkIsValidEmail` - ~250 ns/op
- `BenchmarkValidatePassword` - ~450 ns/op

### ✅ Pagination (`models/pagination_test.go`)

**Функции:**
- `NewPaginationParams()` - 13 тест-кейсов
  - Дефолтные значения, валидные параметры, невалидные (negative, zero, non-numeric)
  - Превышение максимума (cap at MaxLimit)

- `GetOffset()` - 6 тест-кейсов
  - Корректное вычисление offset для SQL

- `NewPaginationResult()` - 9 тест-кейсов
  - Расчет total pages, has_next, has_previous
  - Граничные случаи (пустой результат, одна страница, последняя страница)

**Бенчмарки:**
- `BenchmarkNewPaginationParams` - ~120 ns/op
- `BenchmarkGetOffset` - ~2 ns/op
- `BenchmarkNewPaginationResult` - ~50 ns/op

### ✅ Security Headers (`internal/middleware/security_headers_test.go`)

**Функции:**
- `SecurityHeaders()` - 3 тест-кейса
  - Проверка всех заголовков безопасности
  - HSTS в production vs development
  - CSP директивы

- `buildCSP()` - тестирование Content-Security-Policy

### ✅ Error Handling (`internal/handlers/admin/errors_test.go`)

**Функции:**
- `logAndRespondWithError()` - 4 тест-кейса
  - Internal server error, Not found, Unauthorized, Bad request
  - Проверка что технические детали не утекают клиенту

- `TestErrorConstants` - проверка констант сообщений

---

## 🎓 Как добавить новый тест

### Шаг 1: Создать файл

```bash
# Для файла handler.go создать handler_test.go
touch internal/handlers/admin/your_handler_test.go
```

### Шаг 2: Базовая структура

```go
package admin

import (
    "testing"
)

func TestYourFunction(t *testing.T) {
    tests := []struct {
        name string
        // Другие поля
    }{
        {
            name: "test case 1",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Тест логика
        })
    }
}
```

### Шаг 3: Запустить

```bash
go test -v ./internal/handlers/admin -run TestYourFunction
```

---

## 🔄 CI/CD Integration

### GitHub Actions пример

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Run tests
        run: go test -v -cover ./...

      - name: Check coverage
        run: |
          go test -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out
```

---

## 📚 Дополнительные ресурсы

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Blog: Table Driven Tests](https://go.dev/blog/subtests)
- [Effective Go: Testing](https://go.dev/doc/effective_go#testing)
- [Advanced Testing in Go (video)](https://www.youtube.com/watch?v=8hQG7QlcLBk)

---

## 🐛 Известные проблемы

### Нет integration тестов с БД
В данный момент тесты для репозиториев требуют mock'ов базы данных или тестовую PostgreSQL.

**TODO:**
- [ ] Добавить integration тесты с тестовой БД
- [ ] Настроить docker-compose для тестовой среды
- [ ] Добавить тесты для auth handlers

### Нет E2E тестов
**TODO:**
- [ ] Настроить Selenium/Playwright для E2E
- [ ] Добавить smoke tests для критичных flow

---

## ⚙️ Настройка IDE

### VSCode

`.vscode/settings.json`:
```json
{
  "go.testFlags": ["-v"],
  "go.coverOnSave": true,
  "go.coverageDecorator": {
    "type": "gutter"
  }
}
```

### GoLand

Settings → Go → Testing:
- ✅ Enable coverage on test runs
- ✅ Show coverage in gutter

---

## 📞 Поддержка

При возникновении проблем с тестами:
1. Проверьте версию Go: `go version` (требуется 1.24+)
2. Обновите зависимости: `go mod tidy`
3. Проверьте переменные окружения (если нужны для тестов)
4. Запустите с `-v` для подробного вывода

---

**Последнее обновление:** 2026-01-08
**Версия документа:** 1.0
**Автор:** ProminaYou Team
