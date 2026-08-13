# Загрузка .env
ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: help build run test clean migrate-up migrate-down migrate-create migrate-force migrate-version

# Переменные
BINARY_NAME=script-monitor
BUILD_DIR=bin

# Database URL для миграций
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	go run cmd/server/main.go

# ===== DATABASE =====

db-create: ## Create database if not exists
	@docker exec -i script-monitor-postgres psql -U postgres -c "CREATE DATABASE $(DB_NAME);" 2>/dev/null || true
	@echo "Database ready"

# ===== MIGRATIONS =====

migrate-up: ## Apply all up migrations using psql
	@echo "Applying migrations..."
	@for file in migrations/*.up.sql; do \
		echo "Applying $$file..."; \
		docker exec -i script-monitor-postgres psql -U postgres -d $(DB_NAME) < $$file; \
	done
	@echo "All migrations applied"

migrate-down: ## Rollback one migration
	@echo "Rolling back migration..."
	@migrate -path migrations -database "$(DB_URL)" down 1

migrate-version: ## Show current migration version
	@echo "Current migration version:"
	@migrate -path migrations -database "$(DB_URL)" version

migrate-create: ## Create new migration (usage: make migrate-create NAME=create_table)
	@echo "Creating migration: $(NAME)"
	@migrate create -ext sql -dir migrations -seq $(NAME)

.DEFAULT_GOAL := help