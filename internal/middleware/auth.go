package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

// AuthMiddleware adalah middleware untuk validasi JWT token
func AuthMiddleware(authService services.AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Ambil token dari Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.JSONResponse(w, http.StatusUnauthorized, "Token tidak ditemukan", nil, nil)
				return
			}

			// Format: Bearer <token>
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.JSONResponse(w, http.StatusUnauthorized, "Format token tidak valid", nil, nil)
				return
			}

			tokenString := parts[1]

			// Validasi token dan ambil user ID
			userID, err := authService.GetUserIDFromToken(tokenString)
			if err != nil {
				utils.JSONResponse(w, http.StatusUnauthorized, err.Error(), nil, nil)
				return
			}

			// Ambil role dari token
			role, err := authService.GetUserRoleFromToken(tokenString)
			if err != nil {
				utils.JSONResponse(w, http.StatusUnauthorized, err.Error(), nil, nil)
				return
			}

			// Simpan user ID dan role ke context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// RoleMiddleware adalah middleware untuk validasi role user
func RoleMiddleware(authService services.AuthService, allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Ambil token dari Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.JSONResponse(w, http.StatusUnauthorized, "Token tidak ditemukan", nil, nil)
				return
			}

			// Format: Bearer <token>
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.JSONResponse(w, http.StatusUnauthorized, "Format token tidak valid", nil, nil)
				return
			}

			tokenString := parts[1]

			// Validasi token dan ambil user ID
			userID, err := authService.GetUserIDFromToken(tokenString)
			if err != nil {
				utils.JSONResponse(w, http.StatusUnauthorized, err.Error(), nil, nil)
				return
			}

			// Ambil role dari token
			role, err := authService.GetUserRoleFromToken(tokenString)
			if err != nil {
				utils.JSONResponse(w, http.StatusUnauthorized, err.Error(), nil, nil)
				return
			}

			// Cek apakah role user termasuk dalam allowed roles
			roleAllowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				utils.JSONResponse(w, http.StatusForbidden, "Anda tidak memiliki akses ke endpoint ini", nil, nil)
				return
			}

			// Simpan user ID dan role ke context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// GetUserIDFromContext mengambil user ID dari context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUserRoleFromContext mengambil user role dari context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
