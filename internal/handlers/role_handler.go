package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type RoleHandler struct {
	service services.RoleService
}

func NewRoleHandler(service services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// Route: /api/v1/roles → GET all, POST
func (h *RoleHandler) RolesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetRoles(w, r)
	case http.MethodPost:
		h.CreateRole(w, r)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// Route: /api/v1/roles/{id} → GET, PUT, DELETE
func (h *RoleHandler) RoleByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r, "/api/v1/roles/")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "ID tidak boleh kosong", nil, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetRoleByID(w, r, id)
	case http.MethodPut:
		h.UpdateRole(w, r, id)
	case http.MethodDelete:
		h.DeleteRole(w, r, id)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// GET /api/v1/roles
func (h *RoleHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	meta := &api.Metadata{
		Page:       1,
		Limit:      10,
		Total:      len(roles),
		TotalPages: 1,
	}

	utils.JSONResponse(w, http.StatusOK, "Roles retrieved successfully", roles, meta)
}

// GET /api/v1/roles/{id}
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request, id string) {
	role, err := h.service.GetRoleByID(id)
	if err != nil {
		if err.Error() == "role tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Role retrieved successfully", role, nil)
}

// POST /api/v1/roles
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var input api.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	role, err := h.service.CreateRole(input.Name, input.Description)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "Role created successfully", role, nil)
}

// PUT /api/v1/roles/{id}
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request, id string) {
	var input api.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	role, err := h.service.UpdateRole(id, input.Name, input.Description)
	if err != nil {
		if err.Error() == "role tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Role updated successfully", role, nil)
}

// DELETE /api/v1/roles/{id}
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteRole(id)
	if err != nil {
		if err.Error() == "role tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Role deleted successfully", nil, nil)
}
