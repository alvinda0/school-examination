package handlers

import (
	"net/http"
	"strconv"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type AuditLogHandler struct {
	auditLogService services.AuditLogService
}

func NewAuditLogHandler(auditLogService services.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
	}
}

// GET /api/v1/audit-logs
func (h *AuditLogHandler) GetAllAuditLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("limit"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Build filters from query parameters
	filters := make(map[string]interface{})
	
	if search := query.Get("search"); search != "" {
		filters["search"] = search
	}
	
	if method := query.Get("method"); method != "" {
		filters["method"] = method
	}
	
	if roleID := query.Get("role_id"); roleID != "" {
		filters["role_id"] = roleID
	}
	
	if userID := query.Get("user_id"); userID != "" {
		filters["user_id"] = userID
	}

	logs, total, err := h.auditLogService.GetAllAuditLogs(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response := api.ToAuditLogResponseList(logs)
	meta := &api.Metadata{
		Page:       page,
		Limit:      pageSize,
		Total:      total,
		TotalPages: (total + pageSize - 1) / pageSize,
	}
	utils.JSONResponse(w, http.StatusOK, "Audit logs retrieved successfully", response, meta)
}

// Route: /api/v1/audit-logs → GET
func (h *AuditLogHandler) AuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllAuditLogs(w, r)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}
