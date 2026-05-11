package services

import (
	"errors"
	"strings"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
)

type UserService interface {
	GetAllUsers(roleID string) ([]model.User, error)
	GetUserByID(id string) (*model.User, error)
	CreateUser(fullName, email, password, roleID, status string) (*model.User, error)
	UpdateUser(id, fullName, email, password, roleID, status string) (*model.User, error)
	DeleteUser(id string) error
	UpdateLastLogin(id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
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

func (s *userService) CreateUser(fullName, email, password, roleID, status string) (*model.User, error) {
	if strings.TrimSpace(fullName) == "" {
		return nil, errors.New("full_name tidak boleh kosong")
	}
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email tidak boleh kosong")
	}
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("password tidak boleh kosong")
	}
	if strings.TrimSpace(roleID) == "" {
		return nil, errors.New("role_id tidak boleh kosong")
	}

	// Validasi status
	if status != "" && status != "active" && status != "inactive" && status != "suspended" {
		return nil, errors.New("status harus salah satu dari: active, inactive, suspended")
	}

	// TODO: Hash password sebelum disimpan
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	//     return nil, err
	// }

	return s.repo.Create(fullName, email, password, roleID, status)
}

func (s *userService) UpdateUser(id, fullName, email, password, roleID, status string) (*model.User, error) {
	if strings.TrimSpace(fullName) == "" {
		return nil, errors.New("full_name tidak boleh kosong")
	}
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email tidak boleh kosong")
	}
	if strings.TrimSpace(roleID) == "" {
		return nil, errors.New("role_id tidak boleh kosong")
	}
	if strings.TrimSpace(status) == "" {
		return nil, errors.New("status tidak boleh kosong")
	}

	// Validasi status
	if status != "active" && status != "inactive" && status != "suspended" {
		return nil, errors.New("status harus salah satu dari: active, inactive, suspended")
	}

	// TODO: Hash password jika diubah
	// if password != "" {
	//     hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//     if err != nil {
	//         return nil, err
	//     }
	//     password = string(hashedPassword)
	// }

	user, err := s.repo.Update(id, fullName, email, password, roleID, status)
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
