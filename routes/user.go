package routes

import (
	"school-examination/internal/handlers"
	"school-examination/internal/middleware"
	"school-examination/internal/models"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(api *gin.RouterGroup, h *handlers.UserHandler) {
	// GET /api/v1/me  — semua role yang sudah login
	api.GET("/me", h.GetMe)

	users := api.Group("/users",
		middleware.RequireRoles(models.RoleSuperAdmin, models.RoleAdmin),
	)
	{
		// GET    /api/v1/users          — list semua user (query: ?role=teacher&page=1&limit=20)
		users.GET("", h.GetUsers)

		// POST   /api/v1/users          — buat user baru dengan role apapun (oleh admin)
		// body: { "name", "email", "password", "role": "teacher|student|candidate|admin|super_admin" }
		users.POST("", h.CreateUser)

		// GET    /api/v1/users/:id      — detail user
		users.GET("/:id", h.GetUser)

		// PUT    /api/v1/users/:id      — update nama / role / status aktif
		users.PUT("/:id", h.UpdateUser)

		// DELETE /api/v1/users/:id      — hanya super_admin
		users.DELETE("/:id",
			middleware.RequireRoles(models.RoleSuperAdmin),
			h.DeleteUser,
		)
	}
}
