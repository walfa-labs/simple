# AGENTS.md

This file provides guidance to Qoder (qoder.com) when working with code in this repository.

## Project Overview

A simple Go web application with optional PostgreSQL, Redis, and MongoDB integrations. All database integrations are controlled via `.env` configuration - if credentials are not set, the integration is disabled.

## Commands

```bash
# Build
go build -o bin/app

# Run
go run main.go

# Run with hot reload
air

# Generate Swagger docs
swag init -g main.go -o docs

# Generate SQL code (PostgreSQL)
sqlc generate

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter
golangci-lint run

# Clean build artifacts
make clean
```

## Architecture

```
├── cmd/                  # Application entrypoints
├── config/               # Configuration and .env loading
├── handlers/             # HTTP handlers (Gin)
├── integrations/         # Database integrations
│   ├── postgres/         # PostgreSQL + sqlc generated code
│   ├── redis/            # Redis client
│   └── mongo/            # MongoDB client
├── models/               # Shared data models
├── docs/                 # Swagger documentation (auto-generated)
├── main.go               # Application entrypoint
├── go.mod                # Go module definition
├── Makefile              # Common commands
└── .env.example          # Environment template
```

## Integration Pattern

All integrations follow the same pattern:

1. **Check env vars** → Integration enabled only if required vars exist
2. **Initialize client** → Create connection pool
3. **Register handlers** → Add CRUD endpoints if enabled
4. **Health check** → `/health/{db}` endpoint for each integration

### Environment Variables

```env
# PostgreSQL (optional)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=mydb

# Redis (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# MongoDB (optional)
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_USER=user
MONGO_PASSWORD=pass
MONGO_DB=mydb
```

## Key Libraries

- **Web Framework**: Gin (`github.com/gin-gonic/gin`)
- **Swagger**: swaggo (`github.com/swaggo/swag`, `github.com/swaggo/gin-swagger`, `github.com/swaggo/files`)
- **PostgreSQL**: `github.com/lib/pq`
- **Redis**: `github.com/redis/go-redis/v9`
- **MongoDB**: `go.mongodb.org/mongo-driver`
- **Config**: `github.com/joho/godotenv`

## Development Notes

- All versions in `go.mod` should be latest stable
- Use `air` for hot reload during development (install: `go install github.com/air-verse/air@latest`)
- Swagger docs auto-generated via `swag init -g main.go -o docs`
- Each integration has independent health check endpoint at `/postgres/health`, `/redis/health`, `/mongo/health`
- Access Swagger UI at `http://localhost:8080/swagger/index.html`

## Project Structure

```
simple/
├── config/
│   └── config.go         # Environment configuration loader
├── handlers/
│   └── health.go         # Global health check handler
├── integrations/
│   ├── postgres/
│   │   └── handler.go    # PostgreSQL CRUD + health check
│   ├── redis/
│   │   └── handler.go    # Redis cache operations + health check
│   └── mongo/
│       └── handler.go    # MongoDB CRUD + health check
├── docs/
│   ├── docs.go           # Swagger docs (auto-generated)
│   ├── swagger.json      # Swagger spec (auto-generated)
│   └── swagger.yaml      # Swagger spec (auto-generated)
├── models/               # Shared data models (currently empty)
├── bin/
│   └── app               # Compiled binary (gitignored)
├── main.go               # Application entrypoint
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
├── Makefile              # Common commands
├── .env.example          # Environment template
└── .env                  # Actual environment (gitignored)
```

## Testing

No test files are included. When adding tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestName ./package
```

## Adding New Integrations

To add a new database integration:

1. Create `integrations/<name>/handler.go`
2. Add config struct in `config/config.go`
3. Add initialization in `main.go`
4. Add Swagger comments for new endpoints
5. Re-run `swag init`
