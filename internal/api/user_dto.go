package api

type CreateUserRequest struct {
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	RoleID   *string `json:"role_id,omitempty"` // optional, default: student role
	Status   *bool   `json:"status,omitempty"`  // optional, default: false
}

type PatchUserRequest struct {
	Email  *string `json:"email,omitempty"`
	Status *bool   `json:"status,omitempty"`
}
