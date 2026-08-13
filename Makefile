ifneq (,$(wildcard .env))
    include .env
    export
endif

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: help build run test clean docker-build docker-run docker-dev swagger

BINARY_NAME=script-monitor
BUILD_DIR=bin
DOCKER_IMAGE=script-monitor

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

run:
	@echo "Running $(BINARY_NAME)..."
	go run cmd/server/main.go

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):latest .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker-compose up -d
	@echo "Containers started"

docker-stop:
	@docker-compose down

docker-dev:
	@echo "Running in development mode..."
	@docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

docker-logs:
	@docker-compose logs -f

docker-clean: docker-stop
	@docker-compose down -v
	@docker rmi $(DOCKER_IMAGE):latest 2>/dev/null || true


test:
	go test -race -cover ./...

test-unit:
	go test -race -cover ./test/unit/...

migrate-up:
	@echo "Applying migrations..."
	@migrate -path migrations -database "postgres://postgres:postgres@127.0.0.1:5438/script_monitor?sslmode=disable" up

migrate-down:
	@echo "Current version before down:"
	@migrate -path migrations -database "$(DB_URL)" version
	@echo "Rolling back migration..."
	@migrate -path migrations -database "$(DB_URL)" down 1
	@echo "Version after down:"
	@migrate -path migrations -database "$(DB_URL)" version

migrate-force:
	@echo "Forcing migration version to $(VERSION)..."
	@migrate -path migrations -database "$(DB_URL)" force $(VERSION)

swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs

.DEFAULT_GOAL := help