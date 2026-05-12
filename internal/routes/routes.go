package routes

import (
	"net/http"

	"github.com/alvindashahrul/my-app/internal/handlers"
	"github.com/alvindashahrul/my-app/internal/middleware"
	"github.com/alvindashahrul/my-app/internal/services"
)

func SetupRoutes(userHandler *handlers.UserHandler, roleHandler *handlers.RoleHandler, authHandler *handlers.AuthHandler, studentHandler *handlers.StudentHandler, authService services.AuthService) {
	// Inisialisasi middleware
	authMiddleware := middleware.AuthMiddleware(authService)
	roleMiddleware := func(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
		return middleware.RoleMiddleware(authService, roles...)
	}

	// Auth (login tidak perlu JWT)
	http.HandleFunc("/api/v1/login", authHandler.Login)
	
	// Auth endpoints yang memerlukan JWT (semua role bisa akses)
	http.HandleFunc("/api/v1/auth/me", authMiddleware(authHandler.GetAuthMe))

	// Users CRUD (bisa diakses oleh admin, teacher, student, super_admin)
	http.HandleFunc("/api/v1/users", roleMiddleware("admin", "teacher", "student", "super_admin")(userHandler.UsersHandler))
	http.HandleFunc("/api/v1/users/", roleMiddleware("admin", "teacher", "student", "super_admin")(userHandler.UserByIDHandler))

	// Roles CRUD (bisa diakses oleh admin, teacher, student, super_admin)
	http.HandleFunc("/api/v1/roles", roleMiddleware("admin", "teacher", "student", "super_admin")(roleHandler.RolesHandler))
	http.HandleFunc("/api/v1/roles/", roleMiddleware("admin", "teacher", "student", "super_admin")(roleHandler.RoleByIDHandler))

	// Students (bisa diakses oleh admin, teacher, super_admin)
	http.HandleFunc("/api/v1/students", roleMiddleware("admin", "teacher", "super_admin")(studentHandler.StudentsHandler))
	http.HandleFunc("/api/v1/students/", roleMiddleware("admin", "teacher", "super_admin")(studentHandler.StudentByIDHandler))
}
