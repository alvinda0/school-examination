package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type AuthHandler struct {
	userService services.UserService
	authService services.AuthService
}

func NewAuthHandler(userService services.UserService, authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
	}
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
		return
	}

	var req api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body", nil, nil)
		return
	}

	token, _, _, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		statusCode := http.StatusUnauthorized
		if err.Error() == "akun tidak aktif" {
			statusCode = http.StatusForbidden
		}
		utils.JSONResponse(w, statusCode, err.Error(), nil, nil)
		return
	}

	response := map[string]string{
		"token": token,
	}

	utils.JSONResponse(w, http.StatusOK, "Authentication successful", response, nil)
}

// GET /api/v1/auth/me
func (h *AuthHandler) GetAuthMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
		return
	}

	// Ambil token dari Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.JSONResponse(w, http.StatusUnauthorized, "Token tidak ditemukan", nil, nil)
		return
	}

	// Format: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.JSONResponse(w, http.StatusUnauthorized, "Format token tidak valid", nil, nil)
		return
	}

	tokenString := parts[1]

	// Validasi token dan ambil user ID
	userID, err := h.authService.GetUserIDFromToken(tokenString)
	if err != nil {
		utils.JSONResponse(w, http.StatusUnauthorized, err.Error(), nil, nil)
		return
	}

	userWithRole, err := h.userService.GetUserByIDWithRole(userID)
	if err != nil {
		if err.Error() == "user tidak ditemukan" {
			utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response := api.AuthMeResponse{
		UserID:   userWithRole.UserID,
		FullName: userWithRole.FullName,
		Email:    userWithRole.Email,
		RoleName: userWithRole.RoleName,
		RoleID:   userWithRole.RoleID,
		Status:   userWithRole.Status,
	}

	utils.JSONResponse(w, http.StatusOK, "User info retrieved successfully", response, nil)
}
