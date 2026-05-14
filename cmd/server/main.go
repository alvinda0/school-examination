package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"school-examination/internal/config"
	"school-examination/internal/database"
	"school-examination/internal/handlers"
	"school-examination/internal/repository"
	"school-examination/internal/routes"
	"school-examination/internal/services"
)

func main() {
	// ── Flags ────────────────────────────────────────────────────────
	freshMigrate := flag.Bool("fresh", false, "Reset database: drop semua tabel lalu recreate")
	seedData := flag.Bool("seed", false, "Jalankan seeder setelah migrate")
	flag.Parse()

	// ── Config & DB ──────────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	config.Load()
	db := database.Connect()

	// ── Migration ────────────────────────────────────────────────────
	if *freshMigrate {
		log.Println("⚠️  Fresh migration: dropping and recreating all tables...")
		database.Fresh(db)
	} else {
		database.Migrate(db)
	}

	// ── Seed ─────────────────────────────────────────────────────────
	if *seedData || *freshMigrate {
		database.Seed(db)
	} else {
		database.SeedSuperAdmin(db)
	}

	// ── Repositories ─────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	examRepo := repository.NewExamRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)

	// ── Services ─────────────────────────────────────────────────────
	authService := services.NewAuthService(userRepo)
	examService := services.NewExamService(examRepo, questionRepo, submissionRepo)

	// ── Handlers ─────────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)
	questionHandler := handlers.NewQuestionHandler(questionRepo)
	examHandler := handlers.NewExamHandler(examService, examRepo, submissionRepo)

	// ── Router ───────────────────────────────────────────────────────
	r := routes.Setup(authHandler, userHandler, questionHandler, examHandler)

	// ── Run ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", config.AppConfig.ServerPort)
	log.Printf("Server running on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
