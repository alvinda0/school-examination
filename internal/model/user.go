package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleName adalah tipe string untuk nama role
type RoleName string

const (
	RoleSuperAdmin RoleName = "super_admin"
	RoleAdmin      RoleName = "admin"
	RoleTeacher    RoleName = "teacher"
	RoleStudent    RoleName = "student"
	RoleCandidate  RoleName = "candidate"
)

// Role adalah alias untuk backward-compatibility (middleware, JWT, dll masih pakai model.Role)
type Role = RoleName

var AllRoles = []RoleName{
	RoleSuperAdmin, RoleAdmin, RoleTeacher, RoleStudent, RoleCandidate,
}

// RoleModel adalah tabel roles di database
type RoleModel struct {
	ID          uuid.UUID `json:"id"          gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        RoleName  `json:"name"        gorm:"type:varchar(30);uniqueIndex;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *RoleModel) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// TableName agar GORM pakai tabel "roles"
func (RoleModel) TableName() string {
	return "roles"
}

type User struct {
	ID        uuid.UUID  `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string     `json:"name"       gorm:"not null"`
	Email     string     `json:"email"      gorm:"uniqueIndex;not null"`
	Password  string     `json:"-"          gorm:"not null"`
	RoleID    uuid.UUID  `json:"role_id"    gorm:"type:uuid;not null"`
	RoleModel *RoleModel `json:"role"       gorm:"foreignKey:RoleID"`
	IsActive  bool       `json:"is_active"  gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Role mengembalikan nama role user (helper agar kode lain tidak perlu akses RoleModel langsung)
func (u *User) Role() RoleName {
	if u.RoleModel != nil {
		return u.RoleModel.Name
	}
	return ""
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
	Name     string   `json:"name"     binding:"required"`
	Email    string   `json:"email"    binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Role     RoleName `json:"role"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
