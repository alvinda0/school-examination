// Package appern menyediakan bootstrap dan inisialisasi aplikasi.
// Semua dependency injection dan setup dilakukan di sini.
package appern

import (
	"fmt"
	"log"

	"school-examination/internal/config"
	"school-examination/internal/database"
	"school-examination/internal/handlers"
	"school-examination/internal/repository"
	"school-examination/internal/routes"
	"school-examination/internal/services"

	"github.com/gin-gonic/gin"
)

// App menyimpan semua komponen aplikasi yang sudah diinisialisasi.
type App struct {
	Router *gin.Engine
}

// New menginisialisasi seluruh dependency dan mengembalikan App yang siap dijalankan.
func New() *App {
	// Config & DB
	config.Load()
	db := database.Connect()
	database.Migrate(db)
	database.SeedRoles(db)
	database.SeedSuperAdmin(db)

	// Repositories
	userRepo       := repository.NewUserRepository(db)
	questionRepo   := repository.NewQuestionRepository(db)
	examRepo       := repository.NewExamRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)

	// Services
	authService := services.NewAuthService(userRepo)
	examService := services.NewExamService(examRepo, questionRepo, submissionRepo)

	// Handlers
	authHandler     := handlers.NewAuthHandler(authService)
	userHandler     := handlers.NewUserHandler(userRepo)
	questionHandler := handlers.NewQuestionHandler(questionRepo)
	examHandler     := handlers.NewExamHandler(examService, examRepo, submissionRepo)

	// Router
	r := routes.Setup(authHandler, userHandler, questionHandler, examHandler)

	return &App{Router: r}
}

// Run menjalankan HTTP server.
func (a *App) Run() {
	addr := fmt.Sprintf(":%s", config.AppConfig.ServerPort)
	log.Printf("Server running on http://localhost%s", addr)
	if err := a.Router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
