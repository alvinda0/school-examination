package handlers

import (
	"school-examination/internal/model"
	"school-examination/internal/services"
	"school-examination/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, "Registration successful", resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.OK(c, "Login successful", resp)
}
