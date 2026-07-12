# Bookmark Service

A lightweight Go REST API for creating short-lived bookmark links backed by Redis.

The service exposes health-check, short-link creation, and redirect endpoints. It follows a clean layered architecture so HTTP handling, business logic, repository access, and shared packages stay separated and easy to test.

---

## Overview

This project demonstrates a production-style REST API using:

- Go as the server runtime
- Gin as the HTTP framework
- Redis as the URL storage and dependency checked by health checks
- Environment-based configuration via environment variables or a `.env` file
- Layered architecture: Handler -> Service -> Repository -> Package
- Structured error logging with zerolog
- Unit, package, repository, handler, and integration tests
- Docker and Docker Compose for local development

---

## Features

- Lightweight REST API server
- Health check endpoint with Redis dependency status
- URL shortening endpoint with TTL support
- Redirect endpoint for generated short codes
- Random alphanumeric short-code generation
- Environment-based configuration
- Configurable app port, service name, instance ID, log level, and Redis connection
- Swagger documentation
- Dockerfile and Docker Compose setup
- Makefile shortcuts for common development tasks

---

## Project Structure

```text
bookmark-service/
├── cmd/
│   └── api/
│       └── main.go                  # Application entry point
├── internal/
│   ├── api/                         # API engine and route wiring
│   ├── handler/                     # HTTP handlers
│   │   └── v1/                      # Versioned link handlers and DTOs
│   ├── model/                       # Shared response models
│   ├── repository/                  # Redis-backed persistence adapters
│   ├── service/                     # Business logic
│   └── integration_test/            # Endpoint integration tests
├── pkg/
│   ├── logger/                      # Logging configuration
│   ├── redis/                       # Redis client, config, and test helper
│   └── utils/                       # Shared utilities
├── Dockerfile
├── docker-compose.yaml
├── Makefile
├── go.mod
└── go.sum
```

---

## Requirements

Before running the project, make sure you have:

- Go 1.26 or newer
- Redis
- Docker and Docker Compose, optional but recommended
- Git
- `make`, optional but recommended

---

# Getting Started

## 1. Clone the repository

```bash
git clone https://github.com/khaivutri/bookmark-service.git
cd bookmark-service
```

---

## 2. Install dependencies

Download all Go modules.

```bash
go mod download
```

Or use the Makefile:

```bash
make deps
```

---

## 3. Configure environment variables

You can configure the application in one of two ways:

- Create a `.env` file in the project root.
- Export environment variables directly from your terminal.

Example:

```env
APP_PORT=8080
SERVICE_NAME=bookmark_service
INSTANCE_ID=
LOG_LEVEL=info
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

`INSTANCE_ID` is optional. If omitted or left empty, the application automatically generates a UUID when starting.

When running the app inside Docker Compose, use the Redis service name:

```env
REDIS_ADDRESS=redis:6379
```

---

## 4. Start Redis

Start only Redis with Docker Compose:

```bash
make docker-redis
```

Or run Redis directly with Docker:

```bash
docker run --rm --name redis -p 6379:6379 redis:alpine
```

If a container named `redis` already exists, remove it before starting a new one:

```bash
docker rm -f redis
```

---

## 5. Run the application locally

```bash
make run
```

Override environment values at runtime:

```bash
make run APP_PORT=9090 SERVICE_NAME=bookmark-service-dev LOG_LEVEL=debug
```

Run Swagger generation before starting the app:

```bash
make dev-run
```

---

## 6. Run with Docker Compose

Build and start both Redis and the bookmark service:

```bash
make docker-up
```

Follow logs:

```bash
make docker-logs
```

Stop services:

```bash
make docker-down
```

Equivalent Docker Compose commands:

```bash
docker-compose up --build
docker-compose logs -f
docker-compose down
```

---

## 7. Verify the service

After the server starts successfully, verify the health endpoint:

```bash
curl http://localhost:8080/health-check
```

If you changed `APP_PORT`, replace `8080` with your configured port.

---

# API

## Health Check

| Method | Endpoint        | Description                                  | Success Response |
| ------ | --------------- | -------------------------------------------- | ---------------- |
| GET    | `/health-check` | Returns service health and dependency status | `200 OK`         |

When Redis is unavailable, the endpoint returns `503 Service Unavailable` with dependency status set to `DOWN`.

### Example Request

```http
GET /health-check HTTP/1.1
Host: localhost:8080
```

### Example Healthy Response

```json
{
  "message": "OK",
  "service_name": "bookmark_service",
  "instance_id": "c45f7d4f-f0d0-42dc-90d8-d5eb0f6dbe5e",
  "dependency": {
    "redis": "UP"
  }
}
```

### Example Degraded Response

```json
{
  "message": "DEGRADED",
  "service_name": "bookmark_service",
  "instance_id": "c45f7d4f-f0d0-42dc-90d8-d5eb0f6dbe5e",
  "dependency": {
    "redis": "DOWN"
  }
}
```

---

## Create Short Link

| Method | Endpoint           | Description                          | Success Response |
| ------ | ------------------ | ------------------------------------ | ---------------- |
| POST   | `/v1/links/shorten` | Creates a short code for a given URL | `200 OK`         |

### Request Body

| Field | Type   | Required | Validation               | Description                         |
| ----- | ------ | -------- | ------------------------ | ----------------------------------- |
| `url` | string | Yes      | Must be a valid URL      | Original URL to shorten             |
| `exp` | int64  | Yes      | Must be at least `5`     | Expiration time in seconds          |

### Example Request

```bash
curl -X POST http://localhost:8080/v1/links/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","exp":60}'
```

### Example Response

```json
{
  "code": "AbCDeFK",
  "message": "Shorten URL generated successfully!"
}
```

---

## Redirect

| Method | Endpoint                    | Description                              | Success Response |
| ------ | --------------------------- | ---------------------------------------- | ---------------- |
| GET    | `/v1/links/redirect/{code}`  | Redirects a short code to the stored URL | `302 Found`      |

### Example Request

```bash
curl -i http://localhost:8080/v1/links/redirect/AbCDeFK
```

If the code exists, the response includes a `Location` header pointing to the original URL.

If the code does not exist or has expired, the endpoint returns:

```json
{
  "error": "Code not found"
}
```

---

# Swagger

Generate Swagger docs:

```bash
make swagger
```

Run the app and open:

```text
http://localhost:8080/swagger/index.html
```

---

# Configuration

The application reads configuration from:

- Environment variables
- `.env` file, if available

| Variable         | Default             | Description                                      |
| ---------------- | ------------------- | ------------------------------------------------ |
| `APP_PORT`       | `8080`              | HTTP server port                                 |
| `SERVICE_NAME`   | `bookmark_service`  | Service name returned by the health-check API    |
| `INSTANCE_ID`    | Auto-generated UUID | Optional instance identifier                     |
| `LOG_LEVEL`      | `info`              | Log level passed to zerolog                      |
| `REDIS_ADDRESS`  | `localhost:6379`    | Redis host and port                              |
| `REDIS_PASSWORD` | empty               | Redis password                                   |
| `REDIS_DB`       | `0`                 | Redis database number                            |

`INSTANCE_ID` must be a valid UUID when explicitly provided. Invalid values cause startup to fail.

---

# Makefile Commands

| Command             | Description                                      |
| ------------------- | ------------------------------------------------ |
| `make help`         | Show available targets                           |
| `make deps`         | Download Go modules                              |
| `make tidy`         | Clean up Go module dependencies                  |
| `make test`         | Run tests with coverage report                   |
| `make swagger`      | Generate Swagger docs                            |
| `make run`          | Run the application locally                      |
| `make dev-run`      | Generate Swagger docs, then run the application  |
| `make docker-up`    | Build and start services with Docker Compose     |
| `make docker-down`  | Stop Docker Compose services                     |
| `make docker-logs`  | Follow Docker Compose logs                       |
| `make docker-redis` | Start only Redis with Docker Compose             |
| `make clean`        | Remove build and coverage artifacts              |

---

# Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage through the Makefile:

```bash
make test
```

The test suite includes:

- Handler tests
- Service tests
- Repository tests
- Redis client/config/mock tests
- Logger level tests
- Code generator tests
- Integration tests for health-check and short-link endpoints

---

# License

This project is intended as a starter service for building Go REST APIs and can be extended to fit your own requirements.
