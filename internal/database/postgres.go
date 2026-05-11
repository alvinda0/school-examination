package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/alvindashahrul/my-app/internal/config"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	log.Printf("🔌 Connecting to database: host=%s port=%s dbname=%s user=%s", 
		cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database tidak bisa diakses: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	return db, nil
}
