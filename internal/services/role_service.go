package services

import (
	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/repository"
)

type RoleService interface {
	GetAllRoles() ([]model.Role, error)
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
