package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/middleware"
	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type ClassHandler struct {
	classService   services.ClassService
	studentService services.StudentService
}

func NewClassHandler(classService services.ClassService, studentService services.StudentService) *ClassHandler {
	return &ClassHandler{
		classService:   classService,
		studentService: studentService,
	}
}

func (h *ClassHandler) ClassesHandler(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRoleFromContext(r.Context())

	switch r.Method {
	case http.MethodGet:
		h.GetAllClasses(w, r)
	case http.MethodPost:
		if role == "student" {
			utils.SendErrorResponse(w, http.StatusForbidden, "Anda tidak memiliki akses ke endpoint ini", nil)
			return
		}
		h.CreateClass(w, r)
	default:
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func (h *ClassHandler) ClassByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/classes/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid class ID", nil)
		return
	}

	id, err := uuid.Parse(parts[0])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid class ID format", err)
		return
	}

	role, _ := middleware.GetUserRoleFromContext(r.Context())

	// Handle sub-routes
	if len(parts) > 1 {
		switch parts[1] {
		case "teacher":
			h.GetClassWithTeacher(w, r, id)
		case "students":
			if r.Method == http.MethodGet {
				h.GetClassWithStudents(w, r, id)
			} else if r.Method == http.MethodPost {
				if role == "student" {
					utils.SendErrorResponse(w, http.StatusForbidden, "Anda tidak memiliki akses ke endpoint ini", nil)
					return
				}
				h.AssignStudentsToClass(w, r, id)
			} else {
				utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
			}
		default:
			utils.SendErrorResponse(w, http.StatusNotFound, "Endpoint not found", nil)
		}
		return
	}

	// Handle main class operations
	switch r.Method {
	case http.MethodGet:
		h.GetClassByID(w, r, id)
	case http.MethodPut:
		if role == "student" {
			utils.SendErrorResponse(w, http.StatusForbidden, "Anda tidak memiliki akses ke endpoint ini", nil)
			return
		}
		h.UpdateClass(w, r, id)
	case http.MethodDelete:
		if role == "student" {
			utils.SendErrorResponse(w, http.StatusForbidden, "Anda tidak memiliki akses ke endpoint ini", nil)
			return
		}
		h.DeleteClass(w, r, id)
	default:
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func (h *ClassHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var req api.CreateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	class, err := h.classService.CreateClass(r.Context(), &req)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create class", err)
		return
	}

	r = middleware.SetAuditData(r, &middleware.AuditData{
		Action:     "create",
		EntityID:   &class.ID,
		EntityType: "class",
		NewData: model.JSONB{
			"id":            class.ID,
			"name":          class.Name,
			"grade_level":   class.GradeLevel,
			"academic_year": class.AcademicYear,
			"status":        class.Status,
		},
	})

	utils.SendSuccessResponse(w, http.StatusCreated, "Class created successfully", class)
}

func (h *ClassHandler) GetClassByID(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	class, err := h.classService.GetClassByID(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "Class not found", err)
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Class retrieved successfully", class)
}

func (h *ClassHandler) GetAllClasses(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRoleFromContext(r.Context())

	// Jika student, hanya tampilkan kelas milik student tersebut
	if role == "student" {
		userIDStr, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "User tidak ditemukan", nil)
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "User ID tidak valid", nil)
			return
		}

		student, err := h.studentService.GetStudentByUserID(r.Context(), userID)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusNotFound, "Data student tidak ditemukan", err)
			return
		}
		if student.ClassID == nil {
			utils.SendPaginatedResponse(w, http.StatusOK, "Classes retrieved successfully", []model.ClassWithTeacher{}, 1, 10, 0)
			return
		}

		class, err := h.classService.GetClassByID(r.Context(), *student.ClassID)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusNotFound, "Kelas tidak ditemukan", err)
			return
		}

		classWithTeacher := model.ClassWithTeacher{Class: *class}
		utils.SendPaginatedResponse(w, http.StatusOK, "Classes retrieved successfully", []model.ClassWithTeacher{classWithTeacher}, 1, 10, 1)
		return
	}

	query := r.URL.Query()

	var params api.ClassQueryParams
	params.Page, _ = strconv.Atoi(query.Get("page"))
	params.Limit, _ = strconv.Atoi(query.Get("limit"))

	if gradeLevel := query.Get("grade_level"); gradeLevel != "" {
		if gl, err := strconv.Atoi(gradeLevel); err == nil {
			params.GradeLevel = &gl
		}
	}

	if academicYear := query.Get("academic_year"); academicYear != "" {
		params.AcademicYear = &academicYear
	}

	if status := query.Get("status"); status != "" {
		params.Status = &status
	}

	classes, total, err := h.classService.GetAllClasses(r.Context(), &params)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve classes", err)
		return
	}

	utils.SendPaginatedResponse(w, http.StatusOK, "Classes retrieved successfully", classes, params.Page, params.Limit, total)
}

func (h *ClassHandler) UpdateClass(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var req api.UpdateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Ambil data lama
	oldClass, _ := h.classService.GetClassByID(r.Context(), id)

	class, err := h.classService.UpdateClass(r.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update class", err)
		return
	}

	changes := model.JSONB{}
	if oldClass != nil {
		if req.Name != nil && oldClass.Name != *req.Name {
			changes["name"] = map[string]interface{}{"old": oldClass.Name, "new": *req.Name}
		}
		if req.Status != nil && oldClass.Status != *req.Status {
			changes["status"] = map[string]interface{}{"old": oldClass.Status, "new": *req.Status}
		}
	}
	r = middleware.SetAuditData(r, &middleware.AuditData{
		Action:     "update",
		EntityID:   &id,
		EntityType: "class",
		Changes:    changes,
	})

	utils.SendSuccessResponse(w, http.StatusOK, "Class updated successfully", class)
}

func (h *ClassHandler) DeleteClass(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	// Ambil data sebelum dihapus
	oldClass, _ := h.classService.GetClassByID(r.Context(), id)

	err := h.classService.DeleteClass(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete class", err)
		return
	}

	auditData := &middleware.AuditData{
		Action:     "delete",
		EntityID:   &id,
		EntityType: "class",
	}
	if oldClass != nil {
		auditData.DeletedData = model.JSONB{
			"id":            oldClass.ID,
			"name":          oldClass.Name,
			"grade_level":   oldClass.GradeLevel,
			"academic_year": oldClass.AcademicYear,
		}
	}
	r = middleware.SetAuditData(r, auditData)

	utils.SendSuccessResponse(w, http.StatusOK, "Class deleted successfully", nil)
}

func (h *ClassHandler) GetClassWithTeacher(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	class, err := h.classService.GetClassWithTeacher(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "Class not found", err)
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Class with teacher retrieved successfully", class)
}

func (h *ClassHandler) GetClassWithStudents(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	class, err := h.classService.GetClassWithStudents(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "Class not found", err)
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Class with students retrieved successfully", class)
}

func (h *ClassHandler) AssignStudentsToClass(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var req api.AssignStudentsToClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err := h.classService.AssignStudentsToClass(r.Context(), id, req.StudentIDs)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to assign students to class", err)
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Students assigned to class successfully", nil)
}
