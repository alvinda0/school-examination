package services

import (
	"errors"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
	"github.com/google/uuid"
)

type SubjectService interface {
	GetAllSubjects() ([]model.Subject, error)
	GetAllSubjectsWithTeachers() ([]model.SubjectWithTeachers, error)
	GetSubjectsByClassID(classID uuid.UUID) ([]model.SubjectWithTeachers, error)
	GetSubjectByID(id string) (*model.Subject, error)
	CreateSubject(name string, code, description *string) (*model.Subject, error)
	UpdateSubject(id string, name *string, code, description *string) (*model.Subject, error)
	DeleteSubject(id string) error
}

type subjectService struct {
	repo repository.SubjectRepository
}

func NewSubjectService(repo repository.SubjectRepository) SubjectService {
	return &subjectService{repo: repo}
}

func (s *subjectService) GetAllSubjects() ([]model.Subject, error) {
	return s.repo.GetAll()
}

func (s *subjectService) GetAllSubjectsWithTeachers() ([]model.SubjectWithTeachers, error) {
	return s.repo.GetAllWithTeachers()
}

func (s *subjectService) GetSubjectsByClassID(classID uuid.UUID) ([]model.SubjectWithTeachers, error) {
	return s.repo.GetSubjectsByClassID(classID)
}

func (s *subjectService) GetSubjectByID(id string) (*model.Subject, error) {
	subject, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, errors.New("subject tidak ditemukan")
	}
	return subject, nil
}

func (s *subjectService) CreateSubject(name string, code, description *string) (*model.Subject, error) {
	if name == "" {
		return nil, errors.New("nama subject tidak boleh kosong")
	}

	subject := &model.Subject{
		Name:        name,
		Code:        code,
		Description: description,
	}

	err := s.repo.Create(subject)
	if err != nil {
		return nil, err
	}

	return subject, nil
}

func (s *subjectService) UpdateSubject(id string, name *string, code, description *string) (*model.Subject, error) {
	// Cek apakah subject ada
	subject, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, errors.New("subject tidak ditemukan")
	}

	// Update field yang diberikan
	if name != nil {
		subject.Name = *name
	}
	if code != nil {
		subject.Code = code
	}
	if description != nil {
		subject.Description = description
	}

	err = s.repo.Update(subject)
	if err != nil {
		return nil, err
	}

	return subject, nil
}

func (s *subjectService) DeleteSubject(id string) error {
	// Cek apakah subject ada
	subject, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if subject == nil {
		return errors.New("subject tidak ditemukan")
	}

	return s.repo.Delete(id)
}
