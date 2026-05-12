package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Route: /api/v1/users → GET all, POST
func (h *UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// Route: /api/v1/users/{id} → GET, PATCH, DELETE
func (h *UserHandler) UserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r, "/api/v1/users/")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "ID tidak boleh kosong", nil, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUserByID(w, r, id)
	case http.MethodPatch:
		h.PatchUser(w, r, id)
	case http.MethodDelete:
		h.DeleteUser(w, r, id)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// GET /api/v1/users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	roleID := r.URL.Query().Get("role_id")

	users, err := h.service.GetAllUsers(roleID)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	meta := &api.Metadata{
		Page:       1,
		Limit:      10,
		Total:      len(users),
		TotalPages: 1,
	}

	utils.JSONResponse(w, http.StatusOK, "Users retrieved successfully", users, meta)
}

// GET /api/v1/users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, err := h.service.GetUserByID(id)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "User retrieved successfully", user, nil)
}

// POST /api/v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input api.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	// Default status ke true jika tidak diset
	status := true
	if input.Status != nil {
		status = *input.Status
	}

	user, err := h.service.CreateUser(input.FullName, input.Email, input.Password, input.RoleID, status)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "User created successfully", user, nil)
}

// PATCH /api/v1/users/{id}
func (h *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request, id string) {
	var input api.PatchUserRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Body tidak valid", nil, nil)
		return
	}

	// Validasi minimal satu field harus diisi
	if input.Email == nil && input.Status == nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Minimal satu field (email atau status) harus diisi", nil, nil)
		return
	}

	user, err := h.service.PatchUser(id, input.Email, input.Status)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusBadRequest, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "User updated successfully", user, nil)
}

// DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteUser(id)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	utils.JSONResponse(w, http.StatusOK, "User deleted successfully", nil, nil)
}
