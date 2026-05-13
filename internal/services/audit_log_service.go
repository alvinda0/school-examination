package services

import (
	"context"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
)

type AuditLogService interface {
	GetAllAuditLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.AuditLog, int, error)
	CreateAuditLog(ctx context.Context, log *model.AuditLog) error
}

type auditLogService struct {
	auditLogRepo repository.AuditLogRepository
}

func NewAuditLogService(auditLogRepo repository.AuditLogRepository) AuditLogService {
	return &auditLogService{
		auditLogRepo: auditLogRepo,
	}
}

func (s *auditLogService) GetAllAuditLogs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.AuditLog, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	logs, err := s.auditLogRepo.FindAll(ctx, pageSize, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.auditLogRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (s *auditLogService) CreateAuditLog(ctx context.Context, log *model.AuditLog) error {
	return s.auditLogRepo.Create(ctx, log)
}
