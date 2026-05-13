package api

import (
	"github.com/alvindashahrul/my-app/internal/model"
)

// AuditLogResponse represents the response for an audit log
type AuditLogResponse struct {
	ID          string                 `json:"id"`
	FullName    string                 `json:"full_name"`
	RoleID      string                 `json:"role_id"`
	RoleName    string                 `json:"role_name"`
	Method      string                 `json:"method"`
	Endpoint    string                 `json:"endpoint"`
	StatusCode  int                    `json:"status_code"`
	IPAddress   *string                `json:"ip_address,omitempty"`
	UserAgent   *string                `json:"user_agent,omitempty"`
	DurationMs  *int                   `json:"duration_ms,omitempty"`
	Action      *string                `json:"action,omitempty"`
	EntityID    *string                `json:"entity_id,omitempty"`
	EntityType  *string                `json:"entity_type,omitempty"`
	OldData     map[string]interface{} `json:"old_data,omitempty"`
	NewData     map[string]interface{} `json:"new_data,omitempty"`
	Changes     map[string]interface{} `json:"changes,omitempty"`
	DeletedData map[string]interface{} `json:"deleted_data,omitempty"`
	CreatedAt   string                 `json:"created_at"`
}

// ToAuditLogResponse converts AuditLog model to response DTO
func ToAuditLogResponse(log *model.AuditLog) *AuditLogResponse {
	response := &AuditLogResponse{
		ID:         log.ID.String(),
		FullName:   log.FullName,
		RoleID:     log.RoleID.String(),
		RoleName:   log.RoleName,
		Method:     log.Method,
		Endpoint:   log.Endpoint,
		StatusCode: log.StatusCode,
		IPAddress:  log.IPAddress,
		UserAgent:  log.UserAgent,
		DurationMs: log.DurationMs,
		Action:     log.Action,
		CreatedAt:  log.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if log.EntityID != nil {
		entityIDStr := log.EntityID.String()
		response.EntityID = &entityIDStr
	}

	if log.EntityType != nil {
		response.EntityType = log.EntityType
	}

	if log.OldData != nil {
		response.OldData = log.OldData
	}

	if log.NewData != nil {
		response.NewData = log.NewData
	}

	if log.Changes != nil {
		response.Changes = log.Changes
	}

	if log.DeletedData != nil {
		response.DeletedData = log.DeletedData
	}

	return response
}

// ToAuditLogResponseList converts a slice of AuditLog models to response DTOs
func ToAuditLogResponseList(logs []*model.AuditLog) []*AuditLogResponse {
	responses := make([]*AuditLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = ToAuditLogResponse(log)
	}
	return responses
}
