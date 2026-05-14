package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/appern"
	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
)

type ClassService interface {
	CreateClass(ctx context.Context, req *api.CreateClassRequest) (*model.Class, error)
	GetClassByID(ctx context.Context, id uuid.UUID) (*model.Class, error)
	GetAllClasses(ctx context.Context, params *api.ClassQueryParams) ([]model.ClassWithTeacher, int, error)
	UpdateClass(ctx context.Context, id uuid.UUID, req *api.UpdateClassRequest) (*model.Class, error)
	DeleteClass(ctx context.Context, id uuid.UUID) error
	GetClassWithTeacher(ctx context.Context, id uuid.UUID) (*model.ClassWithTeacherDetail, error)
	GetClassWithStudents(ctx context.Context, id uuid.UUID) (*model.ClassWithStudentsDetail, error)
	AssignStudentsToClass(ctx context.Context, classID uuid.UUID, studentIDs []uuid.UUID) error
	RemoveStudentFromClass(ctx context.Context, studentID uuid.UUID) error
}

type classService struct {
	classRepo   repository.ClassRepository
	studentRepo repository.StudentRepository
	teacherRepo repository.TeacherRepository
}

func NewClassService(
	classRepo repository.ClassRepository,
	studentRepo repository.StudentRepository,
	teacherRepo repository.TeacherRepository,
) ClassService {
	return &classService{
		classRepo:   classRepo,
		studentRepo: studentRepo,
		teacherRepo: teacherRepo,
	}
}

func (s *classService) CreateClass(ctx context.Context, req *api.CreateClassRequest) (*model.Class, error) {
	// Validate homeroom teacher if provided
	if req.HomeroomTeacherID != nil {
		teacher, err := s.teacherRepo.GetByID(req.HomeroomTeacherID.String())
		if err != nil {
			return nil, err
		}
		if teacher == nil {
			return nil, appern.ErrNotFound
		}
	}

	class := &model.Class{
		Name:              req.Name,
		GradeLevel:        req.GradeLevel,
		AcademicYear:      req.AcademicYear,
		HomeroomTeacherID: req.HomeroomTeacherID,
		MaxStudents:       req.MaxStudents,
		Status:            "ACTIVE",
	}

	err := s.classRepo.Create(ctx, class)
	if err != nil {
		return nil, err
	}

	return class, nil
}

func (s *classService) GetClassByID(ctx context.Context, id uuid.UUID) (*model.Class, error) {
	class, err := s.classRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, appern.ErrNotFound
	}
	return class, nil
}

func (s *classService) GetAllClasses(ctx context.Context, params *api.ClassQueryParams) ([]model.ClassWithTeacher, int, error) {
	filters := make(map[string]interface{})

	if params.GradeLevel != nil {
		filters["grade_level"] = *params.GradeLevel
	}
	if params.AcademicYear != nil {
		filters["academic_year"] = *params.AcademicYear
	}
	if params.Status != nil {
		filters["status"] = *params.Status
	}

	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	offset := (params.Page - 1) * params.Limit
	return s.classRepo.GetAllWithTeacher(ctx, filters, params.Limit, offset)
}

func (s *classService) UpdateClass(ctx context.Context, id uuid.UUID, req *api.UpdateClassRequest) (*model.Class, error) {
	// Check if class exists
	class, err := s.classRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, appern.ErrNotFound
	}

	// Validate homeroom teacher if provided
	if req.HomeroomTeacherID != nil {
		teacher, err := s.teacherRepo.GetByID(req.HomeroomTeacherID.String())
		if err != nil {
			return nil, err
		}
		if teacher == nil {
			return nil, errors.New("homeroom teacher not found")
		}
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.GradeLevel != nil {
		updates["grade_level"] = *req.GradeLevel
	}
	if req.AcademicYear != nil {
		updates["academic_year"] = *req.AcademicYear
	}
	if req.HomeroomTeacherID != nil {
		updates["homeroom_teacher_id"] = *req.HomeroomTeacherID
	}
	if req.MaxStudents != nil {
		updates["max_students"] = *req.MaxStudents
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	err = s.classRepo.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	return s.classRepo.GetByID(ctx, id)
}

func (s *classService) DeleteClass(ctx context.Context, id uuid.UUID) error {
	// Check if class exists
	class, err := s.classRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if class == nil {
		return appern.ErrNotFound
	}

	// Check if class has students
	count, err := s.classRepo.GetStudentCount(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete class with students, please remove students first")
	}

	return s.classRepo.Delete(ctx, id)
}

func (s *classService) GetClassWithTeacher(ctx context.Context, id uuid.UUID) (*model.ClassWithTeacherDetail, error) {
	class, err := s.classRepo.GetWithTeacher(ctx, id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, appern.ErrNotFound
	}
	return class, nil
}

func (s *classService) GetClassWithStudents(ctx context.Context, id uuid.UUID) (*model.ClassWithStudentsDetail, error) {
	class, err := s.classRepo.GetWithStudents(ctx, id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, appern.ErrNotFound
	}
	return class, nil
}

func (s *classService) AssignStudentsToClass(ctx context.Context, classID uuid.UUID, studentIDs []uuid.UUID) error {
	// Check if class exists
	class, err := s.classRepo.GetByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return appern.ErrNotFound
	}

	// Check current student count
	currentCount, err := s.classRepo.GetStudentCount(ctx, classID)
	if err != nil {
		return err
	}

	if currentCount+len(studentIDs) > class.MaxStudents {
		return fmt.Errorf("class capacity exceeded: max %d students, current %d, trying to add %d",
			class.MaxStudents, currentCount, len(studentIDs))
	}

	// Assign each student to the class
	for _, studentID := range studentIDs {
		student, err := s.studentRepo.FindByID(ctx, studentID)
		if err != nil {
			return err
		}
		if student == nil {
			return fmt.Errorf("student with ID %s not found", studentID)
		}

		updates := map[string]interface{}{
			"class_id": classID,
		}
		err = s.studentRepo.UpdateFields(ctx, studentID, updates)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *classService) RemoveStudentFromClass(ctx context.Context, studentID uuid.UUID) error {
	student, err := s.studentRepo.FindByID(ctx, studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return appern.ErrNotFound
	}

	updates := map[string]interface{}{
		"class_id": nil,
	}
	return s.studentRepo.UpdateFields(ctx, studentID, updates)
}
