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
	CreateUser(name string, age int) (*model.User, error)
	UpdateUser(id, name string, age int) (*model.User, error)
	DeleteUser(id string) error
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

func (s *userService) CreateUser(name string, age int) (*model.User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name tidak boleh kosong")
	}
	if age <= 0 {
		return nil, errors.New("age harus lebih dari 0")
	}

	return s.repo.Create(name, age)
}

func (s *userService) UpdateUser(id, name string, age int) (*model.User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name tidak boleh kosong")
	}
	if age <= 0 {
		return nil, errors.New("age harus lebih dari 0")
	}

	user, err := s.repo.Update(id, name, age)
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
