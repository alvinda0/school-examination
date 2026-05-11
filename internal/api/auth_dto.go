package api

// AuthMeResponse adalah response untuk endpoint /api/v1/auth/me
type AuthMeResponse struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
	RoleID   string `json:"role_id"`
	Status   bool   `json:"status"`
}

// LoginRequest adalah request untuk login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse adalah response untuk login
type LoginResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	RoleName  string `json:"role_name"`
	RoleID    string `json:"role_id"`
	Status    bool   `json:"status"`
	ExpiresAt int64  `json:"expires_at"`
}
