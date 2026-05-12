package services

import (
	"context"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
	"github.com/google/uuid"
)

type StudentService interface {
	GetStudentByID(ctx context.Context, id uuid.UUID) (*model.Student, error)
	GetAllStudents(ctx context.Context, page, pageSize int) ([]*model.Student, int, error)
}

type studentService struct {
	studentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) StudentService {
	return &studentService{
		studentRepo: studentRepo,
	}
}

func (s *studentService) GetStudentByID(ctx context.Context, id uuid.UUID) (*model.Student, error) {
	return s.studentRepo.FindByID(ctx, id)
}

func (s *studentService) GetAllStudents(ctx context.Context, page, pageSize int) ([]*model.Student, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	students, err := s.studentRepo.FindAll(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.studentRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}
