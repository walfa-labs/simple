# Quickstart: Simple Go API

This guide helps you get started with the Simple Go API codebase.

## Prerequisites

- Go 1.26.1 or later
- (Optional) PostgreSQL, Redis, or MongoDB if you want to use those integrations

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd simple

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env
```

## Configuration

Edit `.env` to configure the integrations you want to use:

### Minimal (no databases)
```env
SERVER_PORT=8080
```

### With PostgreSQL
```env
SERVER_PORT=8080
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=mydb
```

### With Redis
```env
SERVER_PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

### With MongoDB
```env
SERVER_PORT=8080
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_USER=myuser
MONGO_PASSWORD=mypassword
MONGO_DB=mydb
```

## Running the Application

```bash
# Run directly
go run main.go

# Or build and run
go build -o bin/app
./bin/app

# With hot reload (requires air)
air
```

The API will start on `http://localhost:8080` (or the port you configured).

## Accessing the API

### Swagger UI
Open `http://localhost:8080/swagger/index.html` to explore the API interactively.

### Health Check
```bash
curl http://localhost:8080/health
```

### PostgreSQL Example
```bash
# Create a record
curl -X POST http://localhost:8080/postgres/records \
  -H "Content-Type: application/json" \
  -d '{"title": "My First Record"}'

# Get all records
curl http://localhost:8080/postgres/records

# Get specific record
curl http://localhost:8080/postgres/records/1

# Delete record
curl -X DELETE http://localhost:8080/postgres/records/1
```

### Redis Example
```bash
# Set a cache value
curl -X POST http://localhost:8080/redis/cache \
  -H "Content-Type: application/json" \
  -d '{"key": "mykey", "value": "myvalue", "expires_in": 3600}'

# Get cache value
curl http://localhost:8080/redis/cache/mykey

# Delete cache value
curl -X DELETE http://localhost:8080/redis/cache/mykey
```

### MongoDB Example
```bash
# Create a document
curl -X POST http://localhost:8080/mongo/documents \
  -H "Content-Type: application/json" \
  -d '{"title": "My Document"}'

# Get all documents
curl http://localhost:8080/mongo/documents

# Get specific document (replace ID with actual ObjectID)
curl http://localhost:8080/mongo/documents/507f1f77bcf86cd799439011

# Delete document
curl -X DELETE http://localhost:8080/mongo/documents/507f1f77bcf86cd799439011
```

## Development Commands

```bash
# Generate Swagger documentation
swag init -g main.go -o docs

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter
golangci-lint run

# Clean build artifacts
make clean
```

## Project Structure

```
simple/
├── config/              # Configuration and .env loading
├── handlers/            # HTTP handlers (Gin)
├── integrations/        # Database integrations
│   ├── postgres/        # PostgreSQL CRUD operations
│   ├── redis/           # Redis cache operations
│   └── mongo/           # MongoDB CRUD operations
├── docs/                # Swagger documentation (auto-generated)
├── main.go              # Application entrypoint
├── go.mod               # Go module definition
├── Makefile             # Common commands
└── .env.example         # Environment template
```

## Troubleshooting

### Database connection fails
- Check that the database is running
- Verify environment variables in `.env`
- Check the logs for specific error messages

### Port already in use
- Change `SERVER_PORT` in `.env`
- Or kill the process using the port: `lsof -ti:8080 | xargs kill -9`

### Swagger UI not working
- Regenerate docs: `swag init -g main.go -o docs`
- Ensure `_ "simple/docs"` import is present in `main.go`

## Code Review Notes

This codebase was reviewed for quality and the following were noted:

**Strengths**:
- Clean separation of concerns
- Consistent patterns across integrations
- Proper use of context for timeouts
- Good Swagger documentation

**Areas for Improvement**:
- No test coverage
- Missing UPDATE operations for PostgreSQL and MongoDB
- Some unused type definitions
- Error messages may expose internal details

See the full review in `specs/001-codebase-review/`.
