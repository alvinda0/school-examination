package routes

import (
	"net/http"

	"github.com/alvindashahrul/my-app/internal/handlers"
)

func SetupRoutes(userHandler *handlers.UserHandler, roleHandler *handlers.RoleHandler, productHandler *handlers.ProductHandler) {
	// Users CRUD
	http.HandleFunc("/api/v1/users", userHandler.UsersHandler)
	http.HandleFunc("/api/v1/users/", userHandler.UserByIDHandler)

	// Roles CRUD
	http.HandleFunc("/api/v1/roles", roleHandler.RolesHandler)
	http.HandleFunc("/api/v1/roles/", roleHandler.RoleByIDHandler)

	// Products
	http.HandleFunc("/api/v1/products", productHandler.GetProducts)
}
