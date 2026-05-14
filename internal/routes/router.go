package routes

import (
	"log"
	"net/http"
	"strings"
	"time"

	"school-examination/internal/config"
	"school-examination/internal/handlers"
	"school-examination/internal/middleware"
	"school-examination/internal/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(
	authHandler     *handlers.AuthHandler,
	userHandler     *handlers.UserHandler,
	questionHandler *handlers.QuestionHandler,
	examHandler     *handlers.ExamHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// ── CORS ─────────────────────────────────────────────────────────
	origins := strings.Split(config.AppConfig.CORSOrigins, ",")
	log.Printf("CORS allowed origins: %v", origins)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ── Health check ──────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
	})

	// ── Auth (public) ─────────────────────────────────────────────────
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// ── Protected ────────────────────────────────────────────────────
	api := r.Group("/api/v1", middleware.Auth())

	// Me
	api.GET("/me", userHandler.GetMe)

	// Users (super_admin, admin)
	users := api.Group("/users", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin))
	{
		users.GET("", userHandler.GetUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", middleware.RequireRoles(model.RoleSuperAdmin), userHandler.DeleteUser)
	}

	// Subjects
	subjects := api.Group("/subjects")
	{
		subjects.GET("", questionHandler.GetSubjects)
		subjects.POST("", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin), questionHandler.CreateSubject)
	}

	// Questions (super_admin, admin, teacher)
	canManage := middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin, model.RoleTeacher)
	questions := api.Group("/questions", canManage)
	{
		questions.GET("", questionHandler.GetQuestions)
		questions.POST("", questionHandler.CreateQuestion)
		questions.GET("/:id", questionHandler.GetQuestion)
		questions.PUT("/:id", questionHandler.UpdateQuestion)
		questions.DELETE("/:id", questionHandler.DeleteQuestion)
	}

	// Classes
	classes := api.Group("/classes")
	{
		classes.GET("", examHandler.GetClasses)
		classes.POST("", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin), examHandler.CreateClass)
		classes.GET("/:id/students", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin, model.RoleTeacher), examHandler.GetStudentsByClass)
		classes.POST("/assign", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin), examHandler.AssignStudentToClass)
	}

	// Exams (super_admin, admin, teacher)
	exams := api.Group("/exams", canManage)
	{
		exams.POST("", examHandler.CreateExam)
		exams.GET("", examHandler.GetExams)
		exams.GET("/:id", examHandler.GetExam)
		exams.GET("/:id/results", examHandler.GetExamResults)
		exams.POST("/grade-essay", examHandler.GradeEssay)
	}

	// Student
	student := api.Group("/student", middleware.RequireRoles(model.RoleStudent))
	{
		student.GET("/exams", examHandler.GetAvailableExams)
		student.POST("/exams/:id/start", examHandler.StartExam)
		student.POST("/submissions/:submission_id/answer", examHandler.SaveAnswer)
		student.POST("/submissions/:submission_id/submit", examHandler.SubmitExam)
		student.GET("/results", examHandler.GetMyResults)
	}

	// Candidate
	candidate := api.Group("/candidate", middleware.RequireRoles(model.RoleCandidate))
	{
		candidate.GET("/exams", examHandler.GetAvailableExams)
		candidate.POST("/exams/:id/start", examHandler.StartExam)
		candidate.POST("/submissions/:submission_id/answer", examHandler.SaveAnswer)
		candidate.POST("/submissions/:submission_id/submit", examHandler.SubmitExam)
	}

	return r
}
