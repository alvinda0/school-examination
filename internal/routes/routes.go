package routes

import (
	"net/http"

	"github.com/alvindashahrul/my-app/internal/handlers"
	"github.com/alvindashahrul/my-app/internal/middleware"
	"github.com/alvindashahrul/my-app/internal/services"
)

func SetupRoutes(userHandler *handlers.UserHandler, roleHandler *handlers.RoleHandler, authHandler *handlers.AuthHandler, studentHandler *handlers.StudentHandler, subjectHandler *handlers.SubjectHandler, teacherHandler *handlers.TeacherHandler, classHandler *handlers.ClassHandler, auditLogHandler *handlers.AuditLogHandler, authService services.AuthService, auditLogService services.AuditLogService) {
	// Inisialisasi middleware
	authMiddleware := middleware.AuthMiddleware(authService)
	roleMiddleware := func(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
		return middleware.RoleMiddleware(authService, roles...)
	}
	auditMiddleware := middleware.AuditLogMiddleware(authService, auditLogService)

	// Helper: gabungkan role middleware + audit log middleware
	protected := func(handler http.HandlerFunc, roles ...string) http.HandlerFunc {
		return auditMiddleware(roleMiddleware(roles...)(handler))
	}

	// Auth (login tidak perlu JWT)
	http.HandleFunc("/api/v1/login", authHandler.Login)
	
	// Auth endpoints yang memerlukan JWT (semua role bisa akses)
	http.HandleFunc("/api/v1/auth/me", auditMiddleware(authMiddleware(authHandler.GetAuthMe)))

	// Users CRUD (bisa diakses oleh admin, teacher, student, super_admin)
	http.HandleFunc("/api/v1/users", protected(userHandler.UsersHandler, "admin", "teacher", "student", "super_admin"))
	http.HandleFunc("/api/v1/users/", protected(userHandler.UserByIDHandler, "admin", "teacher", "student", "super_admin"))

	// Roles CRUD (bisa diakses oleh admin, teacher, student, super_admin)
	http.HandleFunc("/api/v1/roles", protected(roleHandler.RolesHandler, "admin", "teacher", "student", "super_admin"))
	http.HandleFunc("/api/v1/roles/", protected(roleHandler.RoleByIDHandler, "admin", "teacher", "student", "super_admin"))

	// Students (bisa diakses oleh admin, teacher, super_admin)
	http.HandleFunc("/api/v1/students", protected(studentHandler.StudentsHandler, "admin", "teacher", "super_admin"))
	http.HandleFunc("/api/v1/students/", protected(studentHandler.StudentByIDHandler, "admin", "teacher", "super_admin"))

	// Subjects (bisa diakses oleh admin, teacher, super_admin, dan student untuk GET)
	http.HandleFunc("/api/v1/subjects", protected(subjectHandler.SubjectsHandler, "admin", "teacher", "super_admin", "student"))
	http.HandleFunc("/api/v1/subjects/", protected(subjectHandler.SubjectByIDHandler, "admin", "teacher", "super_admin", "student"))

	// Teachers (bisa diakses oleh admin, super_admin)
	http.HandleFunc("/api/v1/teachers", protected(teacherHandler.TeachersHandler, "admin", "super_admin"))
	http.HandleFunc("/api/v1/teachers/", protected(teacherHandler.TeacherByIDHandler, "admin", "super_admin"))

	// Classes (bisa diakses oleh admin, teacher, super_admin, dan student untuk GET)
	http.HandleFunc("/api/v1/classes", protected(classHandler.ClassesHandler, "admin", "teacher", "super_admin", "student"))
	http.HandleFunc("/api/v1/classes/", protected(classHandler.ClassByIDHandler, "admin", "teacher", "super_admin", "student"))

	// Audit Logs (bisa diakses oleh admin, super_admin)
	http.HandleFunc("/api/v1/audit-logs", protected(auditLogHandler.AuditLogsHandler, "admin", "super_admin"))
}
