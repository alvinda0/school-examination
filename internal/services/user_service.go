package services

import (
	"errors"
	"strings"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers(roleID string) ([]model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetUserByIDWithRole(id string) (*repository.UserWithRole, error)
	CreateUser(fullName, email, password, roleID string, status bool) (*model.User, error)
	PatchUser(id string, email *string, status *bool) (*model.User, error)
	DeleteUser(id string) error
	UpdateLastLogin(id string) error
}

type userService struct {
	repo     repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(repo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userService{
		repo:     repo,
		roleRepo: roleRepo,
	}
}

// hashPassword melakukan hashing password dengan bcrypt
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("gagal melakukan hash password")
	}
	return string(hashedPassword), nil
}

func (s *userService) GetAllUsers(roleID string) ([]model.User, error) {
	return s.repo.GetAll(roleID)
}

func (s *userService) GetUserByID(id string) (*model.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("ID tidak boleh kosong")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return user, nil
}

func (s *userService) CreateUser(fullName, email, password, roleID string, status bool) (*model.User, error) {
	if strings.TrimSpace(fullName) == "" {
		return nil, errors.New("full_name tidak boleh kosong")
	}
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email tidak boleh kosong")
	}
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("password tidak boleh kosong")
	}

	// Jika roleID kosong, gunakan role "candidate" sebagai default
	if strings.TrimSpace(roleID) == "" {
		candidateRole, err := s.roleRepo.GetByName("candidate")
		if err != nil {
			return nil, errors.New("gagal mendapatkan role default")
		}
		if candidateRole == nil {
			return nil, errors.New("role candidate tidak ditemukan")
		}
		roleID = candidateRole.ID
	}

	// Hash password sebelum disimpan
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(fullName, email, hashedPassword, roleID, status)
}

func (s *userService) PatchUser(id string, email *string, status *bool) (*model.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("ID tidak boleh kosong")
	}

	// Validasi email jika diberikan
	if email != nil && strings.TrimSpace(*email) == "" {
		return nil, errors.New("email tidak boleh kosong")
	}

	user, err := s.repo.Patch(id, email, status)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return user, nil
}

func (s *userService) DeleteUser(id string) error {
	rowsAffected, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user tidak ditemukan")
	}

	return nil
}

func (s *userService) UpdateLastLogin(id string) error {
	return s.repo.UpdateLastLogin(id)
}

func (s *userService) GetUserByIDWithRole(id string) (*repository.UserWithRole, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("ID tidak boleh kosong")
	}

	userWithRole, err := s.repo.GetByIDWithRole(id)
	if err != nil {
		return nil, err
	}
	if userWithRole == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return userWithRole, nil
}
