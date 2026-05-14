package services

import (
	"errors"

	"github.com/google/uuid"
	"school-examination/internal/model"
	"school-examination/internal/repository"
	"school-examination/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Cek email sudah ada
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existing.ID != uuid.Nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Default role: student
	role := req.Role
	if role == "" {
		role = model.RoleStudent
	}

	// Self-register hanya boleh sebagai student atau candidate
	if role != model.RoleStudent && role != model.RoleCandidate {
		return nil, errors.New("self-registration only allowed for student or candidate role")
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     role,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}
