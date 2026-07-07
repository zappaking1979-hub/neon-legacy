.PHONY: build run test lint clean dev docker-build docker-up docker-down

APP_NAME = neon-server
BUILD_DIR = ./tmp

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v -count=1

test-race:
	go test ./... -race -count=1

lint:
	go vet ./...
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"

clean:
	rm -rf $(BUILD_DIR)

dev:
	docker compose up -d --build

dev-logs:
	docker compose logs -f

docker-stop:
	docker compose down

docker-reset:
	docker compose down -v

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

coverage:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html

tidy:
	go mod tidy
	go mod verify

all: tidy lint test build
