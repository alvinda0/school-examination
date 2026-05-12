package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type SubjectHandler struct {
	service services.SubjectService
}

func NewSubjectHandler(service services.SubjectService) *SubjectHandler {
	return &SubjectHandler{service: service}
}

// Route: /api/v1/subjects → GET all, POST
func (h *SubjectHandler) SubjectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetSubjects(w, r)
	case http.MethodPost:
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

	switch r.Method {
	case http.MethodGet:
		h.GetSubjectByID(w, r, id)
	case http.MethodPatch:
		h.UpdateSubject(w, r, id)
	case http.MethodDelete:
		h.DeleteSubject(w, r, id)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// GET /api/v1/subjects
func (h *SubjectHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := h.service.GetAllSubjects()
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

	utils.JSONResponse(w, http.StatusCreated, "Subject created successfully", subject, nil)
}

// PATCH /api/v1/subjects/{id}
func (h *SubjectHandler) UpdateSubject(w http.ResponseWriter, r *http.Request, id string) {
	var input api.UpdateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	// Validasi minimal satu field harus diisi
	if input.Name == nil && input.Code == nil && input.Description == nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Minimal satu field harus diisi", nil, nil)
		return
	}

	subject, err := h.service.UpdateSubject(id, input.Name, input.Code, input.Description)
	if err != nil {
		if err.Error() == "subject tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Subject updated successfully", subject, nil)
}

// DELETE /api/v1/subjects/{id}
func (h *SubjectHandler) DeleteSubject(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteSubject(id)
	if err != nil {
		if err.Error() == "subject tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Subject deleted successfully", nil, nil)
}
