package routes

import (
	"school-examination/internal/handlers"
	"school-examination/internal/middleware"
	"school-examination/internal/models"

	"github.com/gin-gonic/gin"
)

func RegisterClassRoutes(api *gin.RouterGroup, h *handlers.ExamHandler) {
	classes := api.Group("/classes")
	{
		// GET  /api/v1/classes                  — semua role yang login
		classes.GET("", h.GetClasses)

		// POST /api/v1/classes                  — super_admin, admin
		classes.POST("",
			middleware.RequireRoles(models.RoleSuperAdmin, models.RoleAdmin),
			h.CreateClass,
		)

		// GET  /api/v1/classes/:id/students     — super_admin, admin, teacher
		classes.GET("/:id/students",
			middleware.RequireRoles(models.RoleSuperAdmin, models.RoleAdmin, models.RoleTeacher),
			h.GetStudentsByClass,
		)

		// POST /api/v1/classes/assign           — super_admin, admin
		// body: { "student_id": 1, "class_id": 2 }
		classes.POST("/assign",
			middleware.RequireRoles(models.RoleSuperAdmin, models.RoleAdmin),
			h.AssignStudentToClass,
		)
	}
}
