.PHONY: build run test clean swagger sqlc-deps deps lint

# Build the application
build:
	go build -o bin/app

# Run the application
run:
	go run main.go

# Run with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Generate Swagger documentation
swagger:
	swag init -g main.go -o docs

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install sqlc dependencies
sqlc-deps:
	@echo "sqlc is a code generator, not a library dependency"
	@echo "Install sqlc: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Full build and test cycle
all: fmt vet lint build test
