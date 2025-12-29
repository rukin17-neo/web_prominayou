# prommsc
prominayou.ru

## Описание проекта
Веб-сайт массажного салона для клиентов.
Так же есть панель администратора - система управления мастерами и услугами для массажного салона. Позволяет администраторам управлять информацией о мастерах, их фотографиями и услугами.

## Технологии
- Go (Golang)
- PostgreSQL
- HTML/CSS/JavaScript
- Bootstrap

## Установка и настройка

### 1. Создание БД.

```sql
-- Создание базы данных
CREATE DATABASE prommsc_db;

-- Подключение к базе данных
\c prommsc_db;

-- Создание пользователя 
CREATE USER prommsc_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE prommsc_db TO prommsc_user;
```

### 2. Создание таблицы услуг (services)

```sql
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    duration TEXT NOT NULL
);
```

### 3. Создание таблицы мастеров (masters)

```sql
CREATE TABLE IF NOT EXISTS masters (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    photo_url TEXT,
    photo_data BYTEA, -- Бинарные данные фотографии
    photo_type TEXT,   -- MIME тип фотографии (image/jpeg, image/png, etc.)
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Основные таблицы:
- `services` - услуги салона
- `masters` - мастера с фотографиями


## Запуск проекта

1. Установите зависимости:
```bash
go mod download
```

2. Создайте и настройте базу данных
 
3. Настройте подключение к базе данных в `config/database.go`

4. Запустите проект:
```bash
go run cmd/web/main.go
```

## Структура проекта

- `cmd/web/` - точка входа в приложение
- `internal/handlers/` - HTTP обработчики
- `models/` - модели данных и репозитории
- `templates/` - HTML шаблоны
- `static/` - статические файлы (CSS, JS, изображения)
- `config/` - конфигурация приложения

## API маршруты

### Клиентская часть:
- `GET /` - главная страница
- `GET /services` - список услуг
- `GET /masters` - список мастеров для клиентов
- `GET /masters/photo/{id}` - получение фотографии мастера
- `GET /contacts` - страница контактов
- `GET /reviews` - страница отзывов

### Административная панель:
- `GET /admin` - главная страница админки
- `GET /admin/masters` - список мастеров
- `POST /admin/masters` - создание/обновление мастера
- `POST /admin/masters/delete` - удаление мастера
- `GET /admin/masters/photo/{id}` - получение фотографии мастера из БД
- `GET /admin/services` - список услуг
- `GET /admin/services/new` - форма создания услуги
- `POST /admin/services` - создание услуги
- `GET /admin/services/edit/{id}` - форма редактирования услуги
- `POST /admin/services/{id}` - обновление услуги
- `POST /admin/services/delete/{id}` - удаление услуги
