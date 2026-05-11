package routes

import (
	"net/http"

	"github.com/alvindashahrul/my-app/internal/handlers"
)

func SetupRoutes(userHandler *handlers.UserHandler, roleHandler *handlers.RoleHandler) {
	// Users CRUD
	http.HandleFunc("/api/v1/users", userHandler.UsersHandler)
	http.HandleFunc("/api/v1/users/", userHandler.UserByIDHandler)

	// Roles CRUD
	http.HandleFunc("/api/v1/roles", roleHandler.RolesHandler)
	http.HandleFunc("/api/v1/roles/", roleHandler.RoleByIDHandler)
}
