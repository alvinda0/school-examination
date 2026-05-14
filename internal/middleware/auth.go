package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"school-examination/internal/model"
	"school-examination/internal/utils"
)

const UserIDKey    = "user_id"
const UserRoleKey  = "user_role"
const UserEmailKey = "user_email"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserRoleKey, string(claims.Role))
		c.Set(UserEmailKey, claims.Email)
		c.Next()
	}
}

func RequireRoles(roles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr, exists := c.Get(UserRoleKey)
		if !exists {
			utils.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}
		userRole := model.Role(roleStr.(string))
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}
		utils.Forbidden(c, "You don't have permission to access this resource")
		c.Abort()
	}
}

func GetUserID(c *gin.Context) uuid.UUID {
	id, _ := c.Get(UserIDKey)
	return id.(uuid.UUID)
}

func GetUserRole(c *gin.Context) model.Role {
	role, _ := c.Get(UserRoleKey)
	return model.Role(role.(string))
}
