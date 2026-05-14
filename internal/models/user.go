package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleTeacher    Role = "teacher"
	RoleStudent    Role = "student"
	RoleCandidate  Role = "candidate"
)

var AllRoles = []Role{
	RoleSuperAdmin, RoleAdmin, RoleTeacher, RoleStudent, RoleCandidate,
}

type User struct {
	ID        uuid.UUID `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name"       gorm:"not null"`
	Email     string    `json:"email"      gorm:"uniqueIndex;not null"`
	Password  string    `json:"-"          gorm:"not null"`
	Role      Role      `json:"role"       gorm:"type:varchar(20);not null;default:'student'"`
	IsActive  bool      `json:"is_active"  gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     Role   `json:"role"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
