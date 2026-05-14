package routes

import (
	"school-examination/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, h *handlers.AuthHandler) {
	auth := r.Group("/api/v1/auth")
	{
		// POST /api/v1/auth/register
		auth.POST("/register", h.Register)

		// POST /api/v1/auth/login
		auth.POST("/login", h.Login)
	}
}
