package model

import "time"

type User struct {
	ID        string     `json:"id"`
	FullName  string     `json:"full_name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"` // tidak ditampilkan di JSON
	RoleID    string     `json:"role_id"`
	Status    string     `json:"status"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
