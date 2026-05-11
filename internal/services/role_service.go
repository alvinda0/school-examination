package services

import (
	"errors"
	"strings"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
)

type RoleService interface {
	GetAllRoles() ([]model.Role, error)
	GetRoleByID(id string) (*model.Role, error)
	CreateRole(name, description string) (*model.Role, error)
	UpdateRole(id, name, description string) (*model.Role, error)
	DeleteRole(id string) error
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) GetAllRoles() ([]model.Role, error) {
	return s.repo.GetAll()
}

func (s *roleService) GetRoleByID(id string) (*model.Role, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("ID tidak boleh kosong")
	}

	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role tidak ditemukan")
	}

	return role, nil
}

func (s *roleService) CreateRole(name, description string) (*model.Role, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name tidak boleh kosong")
	}

	return s.repo.Create(name, description)
}

func (s *roleService) UpdateRole(id, name, description string) (*model.Role, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name tidak boleh kosong")
	}

	role, err := s.repo.Update(id, name, description)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role tidak ditemukan")
	}

	return role, nil
}

func (s *roleService) DeleteRole(id string) error {
	rowsAffected, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("role tidak ditemukan")
	}

	return nil
}
