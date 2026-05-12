package services

import (
	"errors"
	"time"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
	"github.com/google/uuid"
)

type TeacherService interface {
	GetAllTeachers() ([]model.TeacherWithUser, error)
	GetTeacherByID(id string) (*model.TeacherWithSubjects, error)
	CreateTeacher(userID string, nip, gender, birthPlace, religion, phoneNumber, address, photoURL, status *string, birthDate *time.Time) (*model.TeacherWithUser, error)
	UpdateTeacher(id string, nip, gender, birthPlace, religion, phoneNumber, address, photoURL, status *string, birthDate *time.Time) (*model.TeacherWithUser, error)
	DeleteTeacher(id string) error
	AssignSubjects(teacherID string, subjectIDs []string) error
	RemoveSubjects(teacherID string, subjectIDs []string) error
}

type teacherService struct {
	teacherRepo repository.TeacherRepository
	userRepo    repository.UserRepository
	subjectRepo repository.SubjectRepository
}

func NewTeacherService(teacherRepo repository.TeacherRepository, userRepo repository.UserRepository, subjectRepo repository.SubjectRepository) TeacherService {
	return &teacherService{
		teacherRepo: teacherRepo,
		userRepo:    userRepo,
		subjectRepo: subjectRepo,
	}
}

func (s *teacherService) GetAllTeachers() ([]model.TeacherWithUser, error) {
	return s.teacherRepo.GetAll()
}

func (s *teacherService) GetTeacherByID(id string) (*model.TeacherWithSubjects, error) {
	teacher, err := s.teacherRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if teacher == nil {
		return nil, errors.New("teacher tidak ditemukan")
	}
	return teacher, nil
}

func (s *teacherService) CreateTeacher(userID string, nip, gender, birthPlace, religion, phoneNumber, address, photoURL, status *string, birthDate *time.Time) (*model.TeacherWithUser, error) {
	// Validasi user_id
	if userID == "" {
		return nil, errors.New("user_id tidak boleh kosong")
	}

	// Parse UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("user_id tidak valid")
	}

	// Cek apakah user ada
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Cek apakah user sudah menjadi teacher
	existingTeacher, err := s.teacherRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existingTeacher != nil {
		return nil, errors.New("user sudah terdaftar sebagai teacher")
	}

	// Set default status
	defaultStatus := "ACTIVE"
	if status == nil {
		status = &defaultStatus
	}

	teacher := &model.Teacher{
		UserID:      userUUID,
		NIP:         nip,
		Gender:      gender,
		BirthPlace:  birthPlace,
		BirthDate:   birthDate,
		Religion:    religion,
		PhoneNumber: phoneNumber,
		Address:     address,
		PhotoURL:    photoURL,
		Status:      *status,
	}

	err = s.teacherRepo.Create(teacher)
	if err != nil {
		return nil, err
	}

	// Return dengan user data
	return &model.TeacherWithUser{
		Teacher: *teacher,
		User:    *user,
	}, nil
}

func (s *teacherService) UpdateTeacher(id string, nip, gender, birthPlace, religion, phoneNumber, address, photoURL, status *string, birthDate *time.Time) (*model.TeacherWithUser, error) {
	// Cek apakah teacher ada
	teacherData, err := s.teacherRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if teacherData == nil {
		return nil, errors.New("teacher tidak ditemukan")
	}

	teacher := &teacherData.Teacher

	// Update field yang diberikan
	if nip != nil {
		teacher.NIP = nip
	}
	if gender != nil {
		teacher.Gender = gender
	}
	if birthPlace != nil {
		teacher.BirthPlace = birthPlace
	}
	if birthDate != nil {
		teacher.BirthDate = birthDate
	}
	if religion != nil {
		teacher.Religion = religion
	}
	if phoneNumber != nil {
		teacher.PhoneNumber = phoneNumber
	}
	if address != nil {
		teacher.Address = address
	}
	if photoURL != nil {
		teacher.PhotoURL = photoURL
	}
	if status != nil {
		teacher.Status = *status
	}

	err = s.teacherRepo.Update(teacher)
	if err != nil {
		return nil, err
	}

	// Get user data
	user, err := s.userRepo.GetByID(teacher.UserID.String())
	if err != nil {
		return nil, err
	}

	return &model.TeacherWithUser{
		Teacher: *teacher,
		User:    *user,
	}, nil
}

func (s *teacherService) DeleteTeacher(id string) error {
	// Cek apakah teacher ada
	teacher, err := s.teacherRepo.GetByID(id)
	if err != nil {
		return err
	}
	if teacher == nil {
		return errors.New("teacher tidak ditemukan")
	}

	return s.teacherRepo.Delete(id)
}

func (s *teacherService) AssignSubjects(teacherID string, subjectIDs []string) error {
	// Cek apakah teacher ada
	teacher, err := s.teacherRepo.GetByID(teacherID)
	if err != nil {
		return err
	}
	if teacher == nil {
		return errors.New("teacher tidak ditemukan")
	}

	// Validasi semua subject ada
	for _, subjectID := range subjectIDs {
		subject, err := s.subjectRepo.GetByID(subjectID)
		if err != nil {
			return err
		}
		if subject == nil {
			return errors.New("subject dengan id " + subjectID + " tidak ditemukan")
		}
	}

	return s.teacherRepo.AssignSubjects(teacherID, subjectIDs)
}

func (s *teacherService) RemoveSubjects(teacherID string, subjectIDs []string) error {
	// Cek apakah teacher ada
	teacher, err := s.teacherRepo.GetByID(teacherID)
	if err != nil {
		return err
	}
	if teacher == nil {
		return errors.New("teacher tidak ditemukan")
	}

	return s.teacherRepo.RemoveSubjects(teacherID, subjectIDs)
}
