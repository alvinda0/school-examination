package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alvindashahrul/my-app/internal/config"
	"github.com/alvindashahrul/my-app/internal/database"
	"github.com/alvindashahrul/my-app/internal/handlers"
	"github.com/alvindashahrul/my-app/internal/middleware"
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
	subjectRepo := repository.NewSubjectRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	classRepo := repository.NewClassRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, roleRepo)
	roleService := services.NewRoleService(roleRepo)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	studentService := services.NewStudentService(studentRepo)
	subjectService := services.NewSubjectService(subjectRepo)
	teacherService := services.NewTeacherService(teacherRepo, userRepo, subjectRepo)
	classService := services.NewClassService(classRepo, studentRepo, teacherRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	authHandler := handlers.NewAuthHandler(userService, authService)
	studentHandler := handlers.NewStudentHandler(studentService)
	subjectHandler := handlers.NewSubjectHandler(subjectService)
	teacherHandler := handlers.NewTeacherHandler(teacherService)
	classHandler := handlers.NewClassHandler(classService)

	// Setup routes
	routes.SetupRoutes(userHandler, roleHandler, authHandler, studentHandler, subjectHandler, teacherHandler, classHandler, authService)

	// Serve static files (uploaded images)
	fs := http.FileServer(http.Dir("uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// Apply CORS middleware
	corsMiddleware := middleware.CORSMiddleware(cfg.CORSOrigins)
	handler := corsMiddleware(http.DefaultServeMux)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("🚀 Server running at http://localhost%s\n", serverAddr)
	fmt.Printf("🌐 CORS enabled for: %s\n", cfg.CORSOrigins)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
