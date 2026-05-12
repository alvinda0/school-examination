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

# Run migrations using SQL file
migrate:
	@echo "Running migrations..."
	@echo "Note: Make sure PostgreSQL is running on localhost:5433"
	@echo ""
	psql -h localhost -p 5433 -U postgres -d school -f migrations/000001_create_roles.sql
	psql -h localhost -p 5433 -U postgres -d school -f migrations/000001_create_users.sql
	psql -h localhost -p 5433 -U postgres -d school -f migrations/000002_create_students.sql
	psql -h localhost -p 5433 -U postgres -d school -f migrations/000003_insert_sample_students.sql
	@echo ""
	@echo "✅ Migrations completed!"

# Run only student sample data migration
migrate-students:
	@echo "Inserting sample student data..."
	@echo "Note: Make sure PostgreSQL is running on localhost:5433"
	@echo ""
	psql -h localhost -p 5433 -U postgres -d school -f migrations/000003_insert_sample_students.sql
	@echo ""
	@echo "✅ Sample students inserted!"

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
