package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONB type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	FullName    string     `json:"full_name"`
	RoleID      uuid.UUID  `json:"role_id"`
	RoleName    string     `json:"role_name"`
	Method      string     `json:"method"`
	Endpoint    string     `json:"endpoint"`
	StatusCode  int        `json:"status_code"`
	IPAddress   *string    `json:"ip_address,omitempty"`
	UserAgent   *string    `json:"user_agent,omitempty"`
	DurationMs  *int       `json:"duration_ms,omitempty"`
	Action      *string    `json:"action,omitempty"`
	EntityID    *uuid.UUID `json:"entity_id,omitempty"`
	EntityType  *string    `json:"entity_type,omitempty"`
	OldData     JSONB      `json:"old_data,omitempty"`
	NewData     JSONB      `json:"new_data,omitempty"`
	Changes     JSONB      `json:"changes,omitempty"`
	DeletedData JSONB      `json:"deleted_data,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
