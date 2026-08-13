.PHONY: build run clean

BINARY_NAME=server
BUILD_DIR=bin

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

run:
	go run cmd/server/main.go

clean:
	rm -rf $(BUILD_DIR)

help:
	@echo "Available commands:"
	@echo "  make build  - Build the application"
	@echo "  make run    - Build and run the application"
	@echo "  make clean  - Remove build artifacts"
