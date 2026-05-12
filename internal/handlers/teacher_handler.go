package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type TeacherHandler struct {
	service services.TeacherService
}

func NewTeacherHandler(service services.TeacherService) *TeacherHandler {
	return &TeacherHandler{service: service}
}

// Route: /api/v1/teachers → GET all, POST
func (h *TeacherHandler) TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTeachers(w, r)
	case http.MethodPost:
		h.CreateTeacher(w, r)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// Route: /api/v1/teachers/{id} → GET, PATCH, DELETE
func (h *TeacherHandler) TeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/api/v1/teachers/")

	// Check if it's a subjects endpoint
	if strings.Contains(id, "/subjects") {
		teacherID := strings.Split(id, "/")[0]
		h.TeacherSubjectsHandler(w, r, teacherID)
		return
	}

	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "ID tidak boleh kosong", nil, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetTeacherByID(w, r, id)
	case http.MethodPatch:
		h.UpdateTeacher(w, r, id)
	case http.MethodDelete:
		h.DeleteTeacher(w, r, id)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// Route: /api/v1/teachers/{id}/subjects → POST (assign), DELETE (remove)
func (h *TeacherHandler) TeacherSubjectsHandler(w http.ResponseWriter, r *http.Request, teacherID string) {
	switch r.Method {
	case http.MethodPost:
		h.AssignSubjects(w, r, teacherID)
	case http.MethodDelete:
		h.RemoveSubjects(w, r, teacherID)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// GET /api/v1/teachers
func (h *TeacherHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.service.GetAllTeachers()
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	meta := &api.Metadata{
		Page:       1,
		Limit:      100,
		Total:      len(teachers),
		TotalPages: 1,
	}

	utils.JSONResponse(w, http.StatusOK, "Teachers retrieved successfully", teachers, meta)
}

// GET /api/v1/teachers/{id}
func (h *TeacherHandler) GetTeacherByID(w http.ResponseWriter, r *http.Request, id string) {
	teacher, err := h.service.GetTeacherByID(id)
	if err != nil {
		if err.Error() == "teacher tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Teacher retrieved successfully", teacher, nil)
}

// POST /api/v1/teachers
func (h *TeacherHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var input api.CreateTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	if input.UserID == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "user_id tidak boleh kosong", nil, nil)
		return
	}

	teacher, err := h.service.CreateTeacher(
		input.UserID,
		input.NIP,
		input.Gender,
		input.BirthPlace,
		input.Religion,
		input.PhoneNumber,
		input.Address,
		input.PhotoURL,
		input.Status,
		input.BirthDate,
	)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "Teacher created successfully", teacher, nil)
}

// PATCH /api/v1/teachers/{id}
func (h *TeacherHandler) UpdateTeacher(w http.ResponseWriter, r *http.Request, id string) {
	var input api.UpdateTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	teacher, err := h.service.UpdateTeacher(
		id,
		input.NIP,
		input.Gender,
		input.BirthPlace,
		input.Religion,
		input.PhoneNumber,
		input.Address,
		input.PhotoURL,
		input.Status,
		input.BirthDate,
	)
	if err != nil {
		if err.Error() == "teacher tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Teacher updated successfully", teacher, nil)
}

// DELETE /api/v1/teachers/{id}
func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteTeacher(id)
	if err != nil {
		if err.Error() == "teacher tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Teacher deleted successfully", nil, nil)
}

// POST /api/v1/teachers/{id}/subjects
func (h *TeacherHandler) AssignSubjects(w http.ResponseWriter, r *http.Request, teacherID string) {
	var input api.AssignSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	if len(input.SubjectIDs) == 0 {
		utils.JSONResponse(w, http.StatusBadRequest, "subject_ids tidak boleh kosong", nil, nil)
		return
	}

	err := h.service.AssignSubjects(teacherID, input.SubjectIDs)
	if err != nil {
		if err.Error() == "teacher tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Subjects assigned successfully", nil, nil)
}

// DELETE /api/v1/teachers/{id}/subjects
func (h *TeacherHandler) RemoveSubjects(w http.ResponseWriter, r *http.Request, teacherID string) {
	var input api.AssignSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	if len(input.SubjectIDs) == 0 {
		utils.JSONResponse(w, http.StatusBadRequest, "subject_ids tidak boleh kosong", nil, nil)
		return
	}

	err := h.service.RemoveSubjects(teacherID, input.SubjectIDs)
	if err != nil {
		if err.Error() == "teacher tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Subjects removed successfully", nil, nil)
}
