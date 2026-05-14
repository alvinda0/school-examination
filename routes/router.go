package routes

import (
	"net/http"
	"strings"
	"time"

	"school-examination/config"
	"school-examination/internal/handlers"
	"school-examination/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup merakit semua route dan mengembalikan *gin.Engine yang siap dijalankan.
func Setup(
	authHandler     *handlers.AuthHandler,
	userHandler     *handlers.UserHandler,
	questionHandler *handlers.QuestionHandler,
	examHandler     *handlers.ExamHandler,
) *gin.Engine {
	r := gin.Default()

	// ── CORS ────────────────────────────────────────────────────────
	origins := strings.Split(config.AppConfig.CORSOrigins, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ── Health check ─────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
	})

	// ── Public routes (tanpa auth) ───────────────────────────────────
	RegisterAuthRoutes(r, authHandler)

	// ── Protected routes (wajib JWT) ─────────────────────────────────
	api := r.Group("/api/v1", middleware.Auth())

	RegisterUserRoutes(api, userHandler)
	RegisterSubjectRoutes(api, questionHandler)
	RegisterQuestionRoutes(api, questionHandler)
	RegisterClassRoutes(api, examHandler)
	RegisterExamRoutes(api, examHandler)

	return r
}
