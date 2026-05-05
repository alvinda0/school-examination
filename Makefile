.PHONY: run build test clean migrate

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/server.exe cmd/server/main.go

# Run tests
test:
	go test -v ./test/...

# Clean build artifacts
clean:
	rm -rf bin/

# Run migrations (manual)
migrate:
	psql -h localhost -p 5433 -U postgres -d postgres -f migrations/000001_create_users.sql

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Hot reload (requires air)
dev:
	air
