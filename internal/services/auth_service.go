package services

import (
	"errors"
	"strings"
	"time"

	"github.com/alvindashahrul/my-app/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, *repository.UserWithRole, int64, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserIDFromToken(tokenString string) (string, error)
	GetUserRoleFromToken(tokenString string) (string, error)
	GetClaimsFromToken(tokenString string) (map[string]string, error)
	GetFullClaimsFromToken(tokenString string) (map[string]string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		jwtSecret: jwtSecret,
	}
}

// Login melakukan autentikasi user dengan email dan password
func (s *authService) Login(email, password string) (string, *repository.UserWithRole, int64, error) {
	// Validasi input
	if strings.TrimSpace(email) == "" {
		return "", nil, 0, errors.New("email tidak boleh kosong")
	}
	if strings.TrimSpace(password) == "" {
		return "", nil, 0, errors.New("password tidak boleh kosong")
	}

	// Cari user berdasarkan email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, 0, err
	}
	if user == nil {
		return "", nil, 0, errors.New("email atau password salah")
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, 0, errors.New("email atau password salah")
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// Get user with role info
	userWithRole, err := s.userRepo.GetByIDWithRole(user.ID)
	if err != nil {
		return "", nil, 0, err
	}

	// Generate JWT token
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"full_name": userWithRole.FullName,
		"role":      userWithRole.RoleName,
		"role_id":   userWithRole.RoleID,
		"exp":       expiresAt,
		"iat":       time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, 0, errors.New("gagal membuat token")
	}

	return tokenString, userWithRole, expiresAt, nil
}

// ValidateToken memvalidasi JWT token
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return token, nil
}

// GetUserIDFromToken mengambil user ID dari JWT token
func (s *authService) GetUserIDFromToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id tidak ditemukan dalam token")
	}

	return userID, nil
}

// GetUserRoleFromToken mengambil role dari JWT token
func (s *authService) GetUserRoleFromToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role tidak ditemukan dalam token")
	}

	return role, nil
}

// GetClaimsFromToken mengambil semua claims yang dibutuhkan audit log dari JWT token
func (s *authService) GetClaimsFromToken(tokenString string) (map[string]string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	result := map[string]string{
		"user_id":   "",
		"full_name": "",
		"role":      "",
		"role_id":   "",
	}

	if v, ok := claims["user_id"].(string); ok {
		result["user_id"] = v
	}
	if v, ok := claims["full_name"].(string); ok {
		result["full_name"] = v
	}
	if v, ok := claims["role"].(string); ok {
		result["role"] = v
	}
	if v, ok := claims["role_id"].(string); ok {
		result["role_id"] = v
	}

	return result, nil
}

// GetFullClaimsFromToken mengambil claims lengkap, fallback ke DB jika role_id tidak ada di token
func (s *authService) GetFullClaimsFromToken(tokenString string) (map[string]string, error) {
	claims, err := s.GetClaimsFromToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Jika role_id sudah ada di token (token baru), langsung return
	if claims["role_id"] != "" {
		return claims, nil
	}

	// Token lama tidak punya role_id — ambil dari DB pakai user_id
	userWithRole, err := s.userRepo.GetByIDWithRole(claims["user_id"])
	if err != nil || userWithRole == nil {
		return nil, errors.New("user tidak ditemukan")
	}

	claims["role_id"] = userWithRole.RoleID
	claims["full_name"] = userWithRole.FullName
	if claims["role"] == "" {
		claims["role"] = userWithRole.RoleName
	}

	return claims, nil
}
