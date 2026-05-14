package routes

import (
	"school-examination/internal/handlers"
	"school-examination/internal/middleware"
	"school-examination/internal/models"

	"github.com/gin-gonic/gin"
)

func RegisterSubjectRoutes(api *gin.RouterGroup, h *handlers.QuestionHandler) {
	subjects := api.Group("/subjects")
	{
		// GET  /api/v1/subjects   — semua role yang login
		subjects.GET("", h.GetSubjects)

		// POST /api/v1/subjects   — super_admin, admin
		subjects.POST("",
			middleware.RequireRoles(models.RoleSuperAdmin, models.RoleAdmin),
			h.CreateSubject,
		)
	}
}
