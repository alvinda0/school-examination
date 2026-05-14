package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/middleware"
	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
	"github.com/google/uuid"
)

type SubjectHandler struct {
	service        services.SubjectService
	studentService services.StudentService
}

func NewSubjectHandler(service services.SubjectService, studentService services.StudentService) *SubjectHandler {
	return &SubjectHandler{service: service, studentService: studentService}
}

// Route: /api/v1/subjects → GET all, POST
func (h *SubjectHandler) SubjectsHandler(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRoleFromContext(r.Context())

	switch r.Method {
	case http.MethodGet:
		h.GetSubjects(w, r)
	case http.MethodPost:
		if role == "student" {
			utils.JSONResponse(w, http.StatusForbidden, "Anda tidak memiliki akses ke endpoint ini", nil, nil)
			return
		}
		h.CreateSubject(w, r)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// Route: /api/v1/subjects/{id} → GET, PATCH, DELETE
func (h *SubjectHandler) SubjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r, "/api/v1/subjects/")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "ID tidak boleh kosong", nil, nil)
		return
	}

	role, _ := middleware.GetUserRoleFromContext(r.Context())

	switch r.Method {
	case http.MethodGet:
		h.GetSubjectByID(w, r, id)
	case http.MethodPatch:
		if role == "student" {
			utils.JSONResponse(w, http.StatusForbidden, "Anda tidak memiliki akses ke endpoint ini", nil, nil)
			return
		}
		h.UpdateSubject(w, r, id)
	case http.MethodDelete:
		if role == "student" {
			utils.JSONResponse(w, http.StatusForbidden, "Anda tidak memiliki akses ke endpoint ini", nil, nil)
			return
		}
		h.DeleteSubject(w, r, id)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// GET /api/v1/subjects
func (h *SubjectHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRoleFromContext(r.Context())

	// Jika student, hanya tampilkan subjects dari kelas yang dimiliki student
	if role == "student" {
		userIDStr, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok {
			utils.JSONResponse(w, http.StatusUnauthorized, "User tidak ditemukan", nil, nil)
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, "User ID tidak valid", nil, nil)
			return
		}

		student, err := h.studentService.GetStudentByUserID(r.Context(), userID)
		if err != nil {
			utils.JSONResponse(w, http.StatusNotFound, "Data student tidak ditemukan", nil, nil)
			return
		}
		if student.ClassID == nil {
			utils.JSONResponse(w, http.StatusOK, "Subjects retrieved successfully", []model.SubjectWithTeachers{}, &api.Metadata{Page: 1, Limit: 100, Total: 0, TotalPages: 1})
			return
		}

		subjects, err := h.service.GetSubjectsByClassID(*student.ClassID)
		if err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
			return
		}

		meta := &api.Metadata{
			Page:       1,
			Limit:      100,
			Total:      len(subjects),
			TotalPages: 1,
		}
		utils.JSONResponse(w, http.StatusOK, "Subjects retrieved successfully", subjects, meta)
		return
	}

	subjects, err := h.service.GetAllSubjectsWithTeachers()
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	meta := &api.Metadata{
		Page:       1,
		Limit:      100,
		Total:      len(subjects),
		TotalPages: 1,
	}

	utils.JSONResponse(w, http.StatusOK, "Subjects retrieved successfully", subjects, meta)
}

// GET /api/v1/subjects/{id}
func (h *SubjectHandler) GetSubjectByID(w http.ResponseWriter, r *http.Request, id string) {
	subject, err := h.service.GetSubjectByID(id)
	if err != nil {
		if err.Error() == "subject tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Subject retrieved successfully", subject, nil)
}

// POST /api/v1/subjects
func (h *SubjectHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	var input api.CreateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	if input.Name == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "Nama subject tidak boleh kosong", nil, nil)
		return
	}

	subject, err := h.service.CreateSubject(input.Name, input.Code, input.Description)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	entityID := subject.ID
	r = middleware.SetAuditData(r, &middleware.AuditData{
		Action:     "create",
		EntityID:   &entityID,
		EntityType: "subject",
		NewData: model.JSONB{
			"id":   subject.ID,
			"name": subject.Name,
			"code": subject.Code,
		},
	})

	utils.JSONResponse(w, http.StatusCreated, "Subject created successfully", subject, nil)
}

// PATCH /api/v1/subjects/{id}
func (h *SubjectHandler) UpdateSubject(w http.ResponseWriter, r *http.Request, id string) {
	var input api.UpdateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	if input.Name == nil && input.Code == nil && input.Description == nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Minimal satu field harus diisi", nil, nil)
		return
	}

	// Ambil data lama
	oldSubject, _ := h.service.GetSubjectByID(id)

	subject, err := h.service.UpdateSubject(id, input.Name, input.Code, input.Description)
	if err != nil {
		if err.Error() == "subject tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	entityID := subject.ID
	changes := model.JSONB{}
	if oldSubject != nil {
		if input.Name != nil && oldSubject.Name != *input.Name {
			changes["name"] = map[string]interface{}{"old": oldSubject.Name, "new": *input.Name}
		}
		if input.Code != nil {
			changes["code"] = map[string]interface{}{"old": oldSubject.Code, "new": *input.Code}
		}
	}
	r = middleware.SetAuditData(r, &middleware.AuditData{
		Action:     "update",
		EntityID:   &entityID,
		EntityType: "subject",
		Changes:    changes,
	})

	utils.JSONResponse(w, http.StatusOK, "Subject updated successfully", subject, nil)
}

// DELETE /api/v1/subjects/{id}
func (h *SubjectHandler) DeleteSubject(w http.ResponseWriter, r *http.Request, id string) {
	// Ambil data sebelum dihapus
	oldSubject, _ := h.service.GetSubjectByID(id)

	err := h.service.DeleteSubject(id)
	if err != nil {
		if err.Error() == "subject tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	entityID, _ := uuid.Parse(id)
	auditData := &middleware.AuditData{
		Action:     "delete",
		EntityID:   &entityID,
		EntityType: "subject",
	}
	if oldSubject != nil {
		auditData.DeletedData = model.JSONB{
			"id":   oldSubject.ID,
			"name": oldSubject.Name,
			"code": oldSubject.Code,
		}
	}
	r = middleware.SetAuditData(r, auditData)

	utils.JSONResponse(w, http.StatusOK, "Subject deleted successfully", nil, nil)
}
