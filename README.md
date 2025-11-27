# Question-Answer API

RESTful API сервис для вопросов и ответов, построенный на Go с использованием Gin framework, PostgreSQL и Docker.

### Запуск проекта через Docker Compose

```bash
# Клонируйте репозиторий
git clone https://github.com/NKV510/question-answer-api.git
cd question-answer-api

# Запустите приложение
docker compose up -d
# Или используйте
make build
make up
```

Приложение будет доступно по адресу: http://localhost:8080

## API Endpoints

### Questions

- `GET /questions` - Получить все вопросы
- `POST /questions` - Создать новый вопрос
- `GET /questions/:id` - Получить вопрос с ответами
- `DELETE /questions/:id` - Удалить вопрос (с ответами)

### Answers

- `POST /questions/:id/answers` - Добавить ответ к вопросу
- `GET /answers/:id` - Получить конкретный ответ
- `DELETE /answers/:id` - Удалить ответ


## Технологии

- **Go 1.24** - Основной язык программирования
- **Gin** - HTTP web framework
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - База данных
- **Docker** - Контейнеризация
- **Goose** - Миграции базы данных

## Структура проекта

```
question-answer-api/
├── cmd/
│   └── main.go                 # Точка входа приложения
├── internal/
│   ├── config/                 # Конфигурация
│   ├── database/               # Подключение к БД
│   ├── handlers/               # HTTP обработчики
│   ├── models/                 # Модели данных
│   └── repository/             # Слой доступа к данным
├── migrations/                 # Миграции базы данных
├── Dockerfile                  # Конфигурация Docker
└── docker-compose.yml          # Docker Compose
```

## Запуск без Docker

### Требования

- Go 1.24+
- PostgreSQL 15+

### Установка

```bash
# Клонирование репозитория
git clone https://github.com/NKV510/question-answer-api.git
cd question-answer-api

# Установка зависимостей
go mod download

# Настройка базы данных
createdb qa_db

# Запуск миграций
goose -dir migrations postgres "user=postgres dbname=qa_db sslmode=disable" up

# Запуск приложения
go run cmd/main.go
```

## Конфигурация

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=qa_db
DB_SSL_MODE=disable
SERVER_PORT=8080
env=local
```

## Тестирование

```bash
# Запуск unit тестов
go test ./internal/... -v
```

## Примеры запросов

### Создание вопроса

```bash
curl -X POST http://localhost:8080/questions \
  -H "Content-Type: application/json" \
  -d '{"text": "How to learn Go programming?"}'
```

**Response:**
```json
{
  "id": 1,
  "text": "How to learn Go programming?",
  "created_at": "2025-11-27T10:00:00Z"
}
```

### Добавление ответа

```bash
curl -X POST http://localhost:8080/questions/1/answers \
  -H "Content-Type: application/json" \
  -d '{"user_id": "student-123", "text": "Start with the official Go tour!"}'
```

**Response:**
```json
{
  "id": 1,
  "question_id": 1,
  "user_id": "student-123",
  "text": "Start with the official Go tour!",
  "created_at": "2025-11-27T10:05:00Z"
}
```

### Получение вопроса с ответами

```bash
curl http://localhost:8080/questions/1
```

**Response:**
```json
{
  "id": 1,
  "text": "How to learn Go programming?",
  "created_at": "2025-11-27T10:00:00Z",
  "answers": [
    {
      "id": 1,
      "question_id": 1,
      "user_id": "student-123",
      "text": "Start with the official Go tour!",
      "created_at": "2025-11-27T10:05:00Z"
    }
  ]
}
```

## База данных

### Схема данных

```sql
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_answers_question_id ON answers(question_id);
CREATE INDEX idx_answers_user_id ON answers(user_id);
```

### Миграции

Миграции выполняются автоматически при запуске приложения через Goose.

## Устранение неполадок

### Проблемы с подключением к БД

```bash
# Проверка статуса БД
docker compose logs db

# Пересоздание контейнеров
docker compose down
docker compose up -d --build
```

### Проблемы с миграциями

```bash
docker compose exec app goose -dir migrations postgres "user=postgres password=password dbname=qa_db host=db port=5432 sslmode=disable" status
```

## Логирование

Приложение использует структурированное логирование через `slog`:
- В development режиме - текстовый формат
- В production режиме - JSON формат

## Особенности

- Каскадное удаление ответов при удалении вопроса
- Валидация входных данных
- Graceful shutdown

---

**Примечание:** Для работы приложения требуется запущенный экземпляр PostgreSQL. При использовании Docker Compose база данных запускается автоматически.
