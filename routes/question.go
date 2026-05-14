package routes

import (
	"school-examination/internal/handlers"
	"school-examination/internal/middleware"
	"school-examination/internal/models"

	"github.com/gin-gonic/gin"
)

// RegisterQuestionRoutes mendaftarkan route bank soal.
// super_admin & admin bisa kelola semua soal.
// teacher hanya bisa kelola soal yang dia buat sendiri (enforced di handler).
func RegisterQuestionRoutes(api *gin.RouterGroup, h *handlers.QuestionHandler) {
	canManage := middleware.RequireRoles(
		models.RoleSuperAdmin,
		models.RoleAdmin,
		models.RoleTeacher,
	)

	questions := api.Group("/questions", canManage)
	{
		// GET    /api/v1/questions             — list soal (query: ?subject_id=1&page=1&limit=20)
		questions.GET("", h.GetQuestions)

		// POST   /api/v1/questions             — buat soal baru
		questions.POST("", h.CreateQuestion)

		// GET    /api/v1/questions/:id         — detail soal + options
		questions.GET("/:id", h.GetQuestion)

		// PUT    /api/v1/questions/:id         — update soal
		questions.PUT("/:id", h.UpdateQuestion)

		// DELETE /api/v1/questions/:id         — hapus soal
		questions.DELETE("/:id", h.DeleteQuestion)
	}
}
