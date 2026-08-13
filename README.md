# Script Monitor Service

Сервис для удаленного управления bash-скриптами на Linux-хостах через SSH с мониторингом их выполнения через audit subsystem.

## Содержание

- [Архитектура](#архитектура)
- [Технологии](#технологии)
- [Требования](#требования)
- [Быстрый старт](#быстрый-старт)
- [API Endpoints](#api-endpoints)
- [Тестирование](#тестирование)
- [Docker](#docker)
- [Безопасность](#безопасность)
- [Структура проекта](#структура-проекта)

## Архитектура

Проект построен на **чистой архитектуре** с четким разделением на слои:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer (API)                      │
│  ┌─────────────┐  ┌─────────────┐  ┌───────────────────┐  │
│  │ CreateHandler │  │CallbackHandler│  │  HealthHandler   │  │
│  └─────────────┘  └─────────────┘  └───────────────────┘  │
│                           │                                 │
├───────────────────────────┼─────────────────────────────────┤
│                           ▼                                 │
│                    Service Layer                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  ScriptService  │  EventService  │  SSHClient      │   │
│  │  - CreateScript │  - Process     │  - RunCommand   │   │
│  │  - Validate     │    Callback    │  - Retry        │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                 │
├───────────────────────────┼─────────────────────────────────┤
│                           ▼                                 │
│                    Repository Layer                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │      ScriptRepository  │   EventRepository          │   │
│  │  - Create              │   - Create                 │   │
│  │  - GetByPath           │                            │   │
│  │  - Exists              │                            │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                 │
├───────────────────────────┼─────────────────────────────────┤
│                           ▼                                 │
│                    Infrastructure                           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  PostgreSQL  │  SSH  │  Migrations  │  Docker       │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Слои архитектуры

| Слой | Описание |
|------|----------|
| **API (Handler)** | HTTP-обработчики, валидация запросов, формирование ответов |
| **Service** | Бизнес-логика: создание скриптов, обработка событий |
| **Repository** | Интерфейсы для работы с БД, абстракция над хранилищем |
| **PostgreSQL** | Реализация репозиториев для PostgreSQL |
| **SSH Client** | Клиент для удаленного выполнения команд |
| **Domain** | Бизнес-сущности (Script, Event) |

## Технологии

| Компонент | Технология | Версия |
|-----------|------------|--------|
| Язык | Go | 1.26.5 |
| База данных | PostgreSQL | 15 |
| Миграции | golang-migrate | v4.19.1 |
| SSH | crypto/ssh | — |
| Логирование | slog | — |
| Валидация | go-playground/validator | v10 |
| Тестирование | testify | v1.10 |
| Контейнеризация | Docker / Docker Compose | — |

## Быстрый старт

### Локальный запуск

```bash
# 1. Клонируй репозиторий
git clone https://github.com/TelitsynNikita/test-example-for-maksec
cd test-example-for-maksec

# 2. Скопируй конфиг
cp .env.example .env

# 3. Подними PostgreSQL через Docker
docker-compose up -d postgres

# 4. Примени миграции
make migrate-up

# 5. Запусти приложение
make run
```

### Запуск через Docker

```bash
# 1. Собери и запусти
make docker-run

# 2. Проверь логи
make docker-logs

# 3. Останови
make docker-stop
```

### Доступные Makefile команды

```bash
make help          # Показать все доступные команды
make build         # Собрать приложение
make run           # Запустить приложение
make test-unit     # Запустить unit-тесты

# Миграции
make migrate-up    # Применить миграции
make migrate-down  # Откатить миграцию
make migrate-version # Показать текущую версию
make migrate-create NAME=create_table # Создать новую миграцию

# Docker
make docker-build  # Собрать Docker образ
make docker-run    # Запустить Docker контейнеры
make docker-stop   # Остановить Docker контейнеры
make docker-dev    # Запустить в режиме разработки

# Другое
make swagger       # Сгенерировать Swagger документацию
make lint          # Запустить линтер
```

## API Endpoints

### 1. Создание скрипта

**POST** `/create`

Создает bash-скрипт на удаленном хосте.

**Запрос:**

```json
{
  "host": "192.168.1.10",
  "user": "root",
  "password": "password",
  "template": "template1",
  "port": 2222
}
```

| Поле | Тип | Обязательное | Описание |
|------|-----|-------------|----------|
| host | string | ✅ | IP-адрес или hostname удаленного хоста |
| user | string | ✅ | Имя пользователя для SSH |
| password | string | ✅ | Пароль для SSH (мин. 8 символов) |
| template | string | ✅ | Имя шаблона (template1, template2) |
| port | int | ❌ | SSH порт (по умолчанию 22) |

**Ответ (201 Created):**

```json
{
  "script_id": "123e4567-e89b-12d3-a456-426614174000",
  "script_path": "/opt/script-monitor/scripts/123e4567-e89b-12d3-a456-426614174000.sh",
  "created_at": "2026-08-13T00:00:00Z"
}
```

**Ошибки:**

| Код | Описание |
|-----|----------|
| 400 | Невалидный запрос, неверный шаблон |
| 409 | Скрипт с таким путем уже существует |
| 500 | Внутренняя ошибка сервера |

### 2. Callback

**POST** `/callback`

Принимает события от агента мониторинга.

**Запрос:**

```json
{
  "user": "root",
  "script": "/opt/script-monitor/scripts/123e4567-e89b-12d3-a456-426614174000.sh",
  "action": "execute",
  "time": "2026-08-13T00:00:00Z"
}
```

**Заголовки:**

```
Authorization: Bearer <CALLBACK_TOKEN>
```

**Ответ (200 OK):**

```json
{
  "status": "ok"
}
```

**Ошибки:**

| Код | Описание |
|-----|----------|
| 400 | Невалидный запрос |
| 401 | Неверный или отсутствующий токен |
| 404 | Скрипт не найден |
| 500 | Внутренняя ошибка сервера |

### 3. Health Check

**GET** `/health`

Проверка статуса сервиса.

**Ответ (200 OK):**

```json
{
  "status": "ok"
}
```

### 4. Swagger UI

**GET** `/swagger/index.html`

Интерактивная документация API.

## Тестирование

### Unit-тесты

```bash
make test-unit
```

### Тесты с детальным выводом

```bash
go test -v -race -cover ./...
```

## Docker

### Сборка образа

```bash
make docker-build
```

### Запуск

```bash
make docker-run
```

### Остановка

```bash
make docker-stop
```

### Режим разработки (с горячей перезагрузкой)

```bash
make docker-dev
```

### Размер образа

Минимальный размер образа достигается за счет:
- Многоступенчатой сборки (multi-stage)
- Alpine Linux в качестве базового образа
- Статической линковки (CGO_ENABLED=0)
- Удаления отладочной информации (-ldflags="-s -w")

## Безопасность

### Защита от атак

| Механизм | Описание |
|----------|----------|
| **Rate Limiting** | 10 запросов/сек с бурстом 20 |
| **Max Body Size** | Ограничение размера тела запроса (1MB) |
| **Валидация** | Проверка всех входных данных через validator-v10 |
| **DisallowUnknownFields** | Запрет неизвестных полей в JSON |
| **Bulkhead** | Ограничение параллельных SSH подключений (20) |
| **Таймауты** | ReadTimeout 5s, WriteTimeout 10s, SSH Timeout 2s |
| **Bearer Token** | Аутентификация callback эндпоинта |
| **Retry** | Экспоненциальная задержка при ошибках SSH |

### SSH безопасность

- Пароли передаются по запросу, не хранятся в БД
- Поля экранируются перед передачей в shell
- Компенсирующие операции при ошибках

## Структура проекта

```
test-example-for-maksec/
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа
├── internal/
│   ├── api/
│   │   ├── handler/                # HTTP обработчики
│   │   │   ├── create_handler.go
│   │   │   ├── callback_handler.go
│   │   │   └── health_handler.go
│   │   ├── middleware/             # Middleware
│   │   │   ├── auth.go
│   │   │   ├── chain.go
│   │   │   ├── limit.go
│   │   │   └── logging.go
│   │   └── response/               # Форматирование ответов
│   │       └── response.go
│   ├── app/                        # Сборка приложения
│   │   └── app.go
│   ├── config/                     # Конфигурация
│   │   └── config.go
│   ├── domain/                     # Бизнес-сущности
│   │   ├── event.go
│   │   ├── script.go
│   │   ├── template.go
│   │   └── scripts/
│   │       └── templates/          # Шаблоны скриптов (embed)
│   │           ├── template1.sh
│   │           └── template2.sh
│   ├── repository/                 # Репозитории
│   │   ├── postgres/               # PostgreSQL реализация
│   │   │   ├── db.go
│   │   │   ├── event_repo.go
│   │   │   └── script_repo.go
│   │   ├── event_repository.go
│   │   └── script_repository.go
│   ├── server/                     # HTTP сервер
│   │   └── server.go
│   ├── service/                    # Бизнес-логика
│   │   ├── event_service.go
│   │   ├── script_service.go
│   │   └── service.go
│   └── ssh/                        # SSH клиент
│       ├── client.go
│       └── interface.go
├── migrations/                     # Миграции БД
│   ├── 20260813120000_create_scripts_table.up.sql
│   ├── 20260813120000_create_scripts_table.down.sql
│   ├── 20260813120500_create_events_table.up.sql
│   └── 20260813120500_create_events_table.down.sql
├── test/
│   ├── unit/                       # Unit-тесты
│   │   ├── service/
│   │   │   ├── script_service_test.go
│   │   │   └── event_service_test.go
│   │   └── domain/
│   │       └── template_test.go
│   └── integration/                # Интеграционные тесты
│       ├── postgres/
│       │   ├── script_repo_test.go
│       │   └── event_repo_test.go
│       └── setup_test.go
├── docs/                           # Swagger документация
├── .env.example                    # Пример конфигурации
├── .dockerignore
├── .gitignore
├── Dockerfile                      # Docker сборка
├── docker-compose.yml              # Docker Compose
├── docker-compose.dev.yml          # Docker Compose для разработки
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Переменные окружения

```bash
# Server Configuration
SERVER_CREATE_PORT=8080          # Порт для API создания скриптов
SERVER_CALLBACK_PORT=8081        # Порт для callback
SERVER_READ_TIMEOUT=10s          # Таймаут чтения
SERVER_WRITE_TIMEOUT=10s         # Таймаут записи

# Database Configuration
DB_HOST=127.0.0.1               # Хост PostgreSQL
DB_PORT=5438                    # Порт PostgreSQL
DB_USER=postgres                # Пользователь
DB_PASSWORD=postgres            # Пароль
DB_NAME=script_monitor          # Имя БД
DB_SSLMODE=disable              # SSL режим
DB_MAX_CONN=10                  # Максимум соединений

# SSH Configuration
SSH_TIMEOUT=2s                  # Таймаут SSH команд
SSH_PORT=22                     # SSH порт

# Callback Security
CALLBACK_TOKEN=secret           # Токен для callback аутентификации

# Logging
LOG_LEVEL=info                  # Уровень логирования (debug/info/warn/error)
LOG_FORMAT=json                 # Формат (json/text)