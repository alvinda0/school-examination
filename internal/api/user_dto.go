package api

type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
	Status   string `json:"status,omitempty"` // optional, default: active
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	RoleID   string `json:"role_id"`
	Status   string `json:"status"`
}
