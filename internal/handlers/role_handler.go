package handlers

import (
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

// GET /api/v1/role
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
