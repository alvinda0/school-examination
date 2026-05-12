package api

type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
	Status   *bool  `json:"status,omitempty"` // optional, default: true
}

type PatchUserRequest struct {
	Email  *string `json:"email,omitempty"`
	Status *bool   `json:"status,omitempty"`
}
