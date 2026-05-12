package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alvindashahrul/my-app/internal/config"
	"github.com/alvindashahrul/my-app/internal/database"
	"github.com/alvindashahrul/my-app/internal/handlers"
	"github.com/alvindashahrul/my-app/internal/repository"
	"github.com/alvindashahrul/my-app/internal/routes"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found, using default values")
	} else {
		log.Println("✅ .env file loaded")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	studentRepo := repository.NewStudentRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, roleRepo)
	roleService := services.NewRoleService(roleRepo)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	studentService := services.NewStudentService(studentRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	authHandler := handlers.NewAuthHandler(userService, authService)
	studentHandler := handlers.NewStudentHandler(studentService)

	// Setup routes
	routes.SetupRoutes(userHandler, roleHandler, authHandler, studentHandler, authService)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("🚀 Server running at http://localhost%s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
