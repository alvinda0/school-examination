package services

import (
	"context"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
	"github.com/google/uuid"
)

type StudentService interface {
	GetStudentByID(ctx context.Context, id uuid.UUID) (*model.StudentWithUser, error)
	GetAllStudents(ctx context.Context, page, pageSize int) ([]*model.StudentWithUser, int, error)
	CreateStudent(ctx context.Context, student *model.Student) error
	UpdateStudent(ctx context.Context, id uuid.UUID, student *model.Student) error
	DeleteStudent(ctx context.Context, id uuid.UUID) error
}

type studentService struct {
	studentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) StudentService {
	return &studentService{
		studentRepo: studentRepo,
	}
}

func (s *studentService) GetStudentByID(ctx context.Context, id uuid.UUID) (*model.StudentWithUser, error) {
	return s.studentRepo.FindByIDWithUser(ctx, id)
}

func (s *studentService) GetAllStudents(ctx context.Context, page, pageSize int) ([]*model.StudentWithUser, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	students, err := s.studentRepo.FindAllWithUser(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.studentRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

func (s *studentService) CreateStudent(ctx context.Context, student *model.Student) error {
	return s.studentRepo.Create(ctx, student)
}

func (s *studentService) UpdateStudent(ctx context.Context, id uuid.UUID, student *model.Student) error {
	// Verify student exists
	_, err := s.studentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	student.ID = id
	return s.studentRepo.Update(ctx, student)
}

func (s *studentService) DeleteStudent(ctx context.Context, id uuid.UUID) error {
	// Verify student exists
	_, err := s.studentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.studentRepo.Delete(ctx, id)
}