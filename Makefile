.PHONY: run build migrate migrate-fresh migrate-seed migrate-drop test tidy

# Jalankan server (development)
run:
	go run main.go

# Build binary
build:
	go build -o bin/server main.go

# Migrate aman (tidak hapus data)
migrate:
	go run cmd/migrate/main.go

# Fresh migration (HAPUS SEMUA DATA)
migrate-fresh:
	go run cmd/migrate/main.go --fresh

# Migrate + seed data awal
migrate-seed:
	go run cmd/migrate/main.go --seed

# Fresh + seed
migrate-fresh-seed:
	go run cmd/migrate/main.go --fresh --seed

# Drop semua tabel
migrate-drop:
	go run cmd/migrate/main.go --drop

# Jalankan test
test:
	go test ./...

# Tidy dependencies
tidy:
	go mod tidy

# Hot reload dengan air
dev:
	air
