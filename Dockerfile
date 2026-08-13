# Stage 1: Builder
FROM golang:1.26.5-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git make

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/bin/server ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.19 AS runner

# Устанавливаем ca-certificates для HTTPS и tzdata для временных зон
RUN apk add --no-cache ca-certificates tzdata

# Создаем пользователя для запуска приложения (без прав root)
RUN adduser -D -u 1000 appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/bin/server /app/server

# Копируем миграции (если нужно)
COPY --from=builder /app/migrations /app/migrations

# Меняем владельца на appuser
RUN chown -R appuser:appuser /app

# Переключаемся на пользователя appuser
USER appuser

# Проксируем порты
EXPOSE 8080 8081

# Запускаем приложение
CMD ["/app/server"]