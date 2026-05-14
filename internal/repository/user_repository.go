package repository

import (
	"errors"

	"school-examination/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// withRole preload RoleModel agar user.Role() bisa dipakai
func (r *UserRepository) withRole() *gorm.DB {
	return r.db.Preload("RoleModel")
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.withRole().Where("users.email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.withRole().First(&user, "users.id = ?", id).Error
	return &user, err
}

func (r *UserRepository) FindAll(page, limit int, roleName string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})

	if roleName != "" {
		// join ke tabel roles untuk filter by nama role
		query = query.
			Joins("JOIN roles ON roles.id = users.role_id").
			Where("roles.name = ?", roleName)
	}

	query.Count(&total)
	err := query.Preload("RoleModel").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

// FindRoleByName mencari role berdasarkan nama
func (r *UserRepository) FindRoleByName(name model.RoleName) (*model.RoleModel, error) {
	var role model.RoleModel
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, errors.New("role not found: " + string(name))
	}
	return &role, nil
}
