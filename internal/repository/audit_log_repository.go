package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alvindashahrul/my-app/internal/model"
)

type AuditLogRepository interface {
	FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*model.AuditLog, error)
	Count(ctx context.Context, filters map[string]interface{}) (int, error)
	Create(ctx context.Context, log *model.AuditLog) error
}

type auditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*model.AuditLog, error) {
	query := `
		SELECT 
			id, user_id, full_name, role_id, role_name, method, endpoint, 
			status_code, ip_address, user_agent, duration_ms, action, 
			entity_id, entity_type, old_data, new_data, changes, deleted_data, 
			created_at
		FROM audit_logs
		WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	// Filter by search (full_name or endpoint)
	if search, ok := filters["search"].(string); ok && search != "" {
		query += ` AND (full_name ILIKE $` + fmt.Sprintf("%d", argPos) + ` OR endpoint ILIKE $` + fmt.Sprintf("%d", argPos) + `)`
		args = append(args, "%"+search+"%")
		argPos++
	}

	// Filter by method
	if method, ok := filters["method"].(string); ok && method != "" {
		query += ` AND method = $` + fmt.Sprintf("%d", argPos)
		args = append(args, method)
		argPos++
	}

	// Filter by role_id
	if roleID, ok := filters["role_id"].(string); ok && roleID != "" {
		query += ` AND role_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, roleID)
		argPos++
	}

	// Filter by user_id
	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		query += ` AND user_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, userID)
		argPos++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argPos) + ` OFFSET $` + fmt.Sprintf("%d", argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*model.AuditLog
	for rows.Next() {
		var log model.AuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.FullName, &log.RoleID, &log.RoleName,
			&log.Method, &log.Endpoint, &log.StatusCode, &log.IPAddress,
			&log.UserAgent, &log.DurationMs, &log.Action, &log.EntityID,
			&log.EntityType, &log.OldData, &log.NewData, &log.Changes,
			&log.DeletedData, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

func (r *auditLogRepository) Count(ctx context.Context, filters map[string]interface{}) (int, error) {
	query := `SELECT COUNT(*) FROM audit_logs WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	// Filter by search (full_name or endpoint)
	if search, ok := filters["search"].(string); ok && search != "" {
		query += ` AND (full_name ILIKE $` + fmt.Sprintf("%d", argPos) + ` OR endpoint ILIKE $` + fmt.Sprintf("%d", argPos) + `)`
		args = append(args, "%"+search+"%")
		argPos++
	}

	// Filter by method
	if method, ok := filters["method"].(string); ok && method != "" {
		query += ` AND method = $` + fmt.Sprintf("%d", argPos)
		args = append(args, method)
		argPos++
	}

	// Filter by role_id
	if roleID, ok := filters["role_id"].(string); ok && roleID != "" {
		query += ` AND role_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, roleID)
		argPos++
	}

	// Filter by user_id
	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		query += ` AND user_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, userID)
		argPos++
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *auditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			id, user_id, full_name, role_id, role_name, method, endpoint,
			status_code, ip_address, user_agent, duration_ms, action,
			entity_id, entity_type, old_data, new_data, changes, deleted_data,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.UserID, log.FullName, log.RoleID, log.RoleName,
		log.Method, log.Endpoint, log.StatusCode, log.IPAddress,
		log.UserAgent, log.DurationMs, log.Action, log.EntityID,
		log.EntityType, log.OldData, log.NewData, log.Changes,
		log.DeletedData, log.CreatedAt,
	)
	return err
}
