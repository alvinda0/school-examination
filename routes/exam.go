package routes

import (
	"school-examination/internal/handlers"
	"school-examination/internal/middleware"
	"school-examination/internal/models"

	"github.com/gin-gonic/gin"
)

// RegisterExamRoutes mendaftarkan route ujian untuk guru/admin,
// route siswa mengerjakan ujian, dan route candidate untuk ujian seleksi.
func RegisterExamRoutes(api *gin.RouterGroup, h *handlers.ExamHandler) {
	canManage := middleware.RequireRoles(
		models.RoleSuperAdmin,
		models.RoleAdmin,
		models.RoleTeacher,
	)

	// ── Manajemen Ujian (guru / admin) ──────────────────────────────
	exams := api.Group("/exams", canManage)
	{
		// POST /api/v1/exams                    — buat jadwal ujian baru
		exams.POST("", h.CreateExam)

		// GET  /api/v1/exams                    — list ujian (query: ?subject_id=&class_id=&page=&limit=)
		exams.GET("", h.GetExams)

		// GET  /api/v1/exams/:id                — detail ujian + soal
		exams.GET("/:id", h.GetExam)

		// GET  /api/v1/exams/:id/results        — rekap nilai semua siswa
		exams.GET("/:id/results", h.GetExamResults)

		// POST /api/v1/exams/grade-essay        — nilai jawaban esai manual
		// body: { "answer_id": 1, "score": 80 }
		exams.POST("/grade-essay", h.GradeEssay)
	}

	// ── Ujian Siswa ──────────────────────────────────────────────────
	// student: siswa terdaftar, bisa lihat nilai sendiri
	student := api.Group("/student",
		middleware.RequireRoles(models.RoleStudent),
	)
	{
		// GET  /api/v1/student/exams                              — ujian yang tersedia
		student.GET("/exams", h.GetAvailableExams)

		// POST /api/v1/student/exams/:id/start                   — mulai sesi ujian
		student.POST("/exams/:id/start", h.StartExam)

		// POST /api/v1/student/submissions/:submission_id/answer  — simpan jawaban (auto-save)
		// body: { "question_id": 1, "selected_option": 3 }
		student.POST("/submissions/:submission_id/answer", h.SaveAnswer)

		// POST /api/v1/student/submissions/:submission_id/submit  — submit ujian
		// body: { "answers": [...] }
		student.POST("/submissions/:submission_id/submit", h.SubmitExam)

		// GET  /api/v1/student/results                            — riwayat & nilai ujian sendiri
		student.GET("/results", h.GetMyResults)
	}

	// ── Ujian Candidate ──────────────────────────────────────────────
	// candidate: calon siswa, hanya bisa ikut ujian seleksi/penerimaan
	// tidak bisa lihat nilai (hasil diumumkan oleh admin)
	candidate := api.Group("/candidate",
		middleware.RequireRoles(models.RoleCandidate),
	)
	{
		// GET  /api/v1/candidate/exams                            — ujian seleksi yang tersedia
		candidate.GET("/exams", h.GetAvailableExams)

		// POST /api/v1/candidate/exams/:id/start                 — mulai ujian seleksi
		candidate.POST("/exams/:id/start", h.StartExam)

		// POST /api/v1/candidate/submissions/:submission_id/answer — simpan jawaban (auto-save)
		candidate.POST("/submissions/:submission_id/answer", h.SaveAnswer)

		// POST /api/v1/candidate/submissions/:submission_id/submit — submit ujian seleksi
		candidate.POST("/submissions/:submission_id/submit", h.SubmitExam)
	}
}
