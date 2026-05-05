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
)

func main() {
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

	// Initialize services
	userService := services.NewUserService(userRepo)
	roleService := services.NewRoleService(roleRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)

	// Setup routes
	routes.SetupRoutes(userHandler, roleHandler)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Server jalan di http://localhost%s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
