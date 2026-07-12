.PHONY: help deps tidy test run dev-run clean swagger docker-up docker-down docker-logs docker-redis

GO ?= go
DOCKER_COMPOSE ?= docker-compose
MAIN_PATH ?= ./cmd/api
BIN_DIR ?= ./bin
BINARY ?= $(BIN_DIR)/bookmark-service
APP_PORT ?= 8080
SERVICE_NAME ?= bookmark_service
INSTANCE_ID ?=
LOG_LEVEL ?= info

help:
	@echo "Available targets:"
	@echo "  make help              Show this help message"
	@echo "  make deps              Download Go modules"
	@echo "  make tidy              Clean up module dependencies"
	@echo "  make test     			Run tests with coverage report"
	@echo "  make run               Run the application"
	@echo "  make clean             Remove build artifacts and coverage files"
	@echo "  make dev-run           Run swagger then run the application (development)"
	@echo "  make docker-up         Build and start services with docker-compose"
	@echo "  make docker-down       Stop docker-compose services"
	@echo "  make docker-logs       Follow docker-compose logs"
	@echo "  make docker-redis      Start only redis with docker-compose"

deps:
	$(GO) mod download

tidy:
	$(GO) mod tidy


COVERAGE_EXCLUDE=mocks|main|test|docs
test:
	$(GO) test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

swagger:
	swag init -g cmd/api/main.go --output docs

# run the application. Only pass INSTANCE_ID into the environment if it was
# explicitly provided to avoid overriding values from a .env file.
run:
	@echo "Starting application (APP_PORT=$(APP_PORT) SERVICE_NAME=$(SERVICE_NAME) INSTANCE_ID=$(INSTANCE_ID) LOG_LEVEL=$(LOG_LEVEL))"
	@if [ -z "$(INSTANCE_ID)" ]; then \
		APP_PORT=$(APP_PORT) SERVICE_NAME=$(SERVICE_NAME) LOG_LEVEL=$(LOG_LEVEL) $(GO) run $(MAIN_PATH); \
	else \
		APP_PORT=$(APP_PORT) SERVICE_NAME=$(SERVICE_NAME) INSTANCE_ID=$(INSTANCE_ID) LOG_LEVEL=$(LOG_LEVEL) $(GO) run $(MAIN_PATH); \
	fi

dev-run: swagger run

docker-up:
	$(DOCKER_COMPOSE) up --build

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-redis:
	$(DOCKER_COMPOSE) up redis

clean:
	rm -rf $(BIN_DIR) coverage.out coverage.html coverage.tmp




