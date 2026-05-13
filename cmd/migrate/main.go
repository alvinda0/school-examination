package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alvindashahrul/my-app/internal/config"
	"github.com/alvindashahrul/my-app/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get migration file from argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go <migration-file>")
	}

	migrationFile := os.Args[1]
	
	// Read migration file
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	// Execute migration
	fmt.Printf("🔄 Running migration: %s\n", filepath.Base(migrationFile))
	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Println("✅ Migration completed successfully!")
}
