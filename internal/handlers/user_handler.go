package handlers

import (
	"strconv"

	"school-examination/internal/middleware"
	"school-examination/internal/models"
	"school-examination/internal/repository"
	"school-examination/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}
	utils.OK(c, "Profile fetched", user)
}

// CreateUser — admin membuat user dengan role apapun
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if req.Role == "" {
		req.Role = models.RoleStudent
	}
	if !isValidRole(req.Role) {
		utils.BadRequest(c, "Invalid role. Valid: super_admin, admin, teacher, student, candidate")
		return
	}

	existing, err := h.userRepo.FindByEmail(req.Email)
	if err == nil && existing.ID != uuid.Nil {
		utils.BadRequest(c, "Email already registered")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
		IsActive: true,
	}
	if err := h.userRepo.Create(user); err != nil {
		utils.InternalError(c, "Failed to create user")
		return
	}
	utils.Created(c, "User created", user)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := parseInt(c.DefaultQuery("page", "1"))
	limit, _ := parseInt(c.DefaultQuery("limit", "20"))
	role := c.Query("role")

	users, total, err := h.userRepo.FindAll(page, limit, role)
	if err != nil {
		utils.InternalError(c, "Failed to fetch users")
		return
	}
	utils.Paginated(c, "Users fetched", users, total, page, limit)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}
	utils.OK(c, "User fetched", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}
	user, err := h.userRepo.FindByID(id)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	var body struct {
		Name     string      `json:"name"`
		Role     models.Role `json:"role"`
		IsActive *bool       `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if body.Name != "" {
		user.Name = body.Name
	}
	if body.IsActive != nil {
		user.IsActive = *body.IsActive
	}
	if body.Role != "" {
		if !isValidRole(body.Role) {
			utils.BadRequest(c, "Invalid role. Valid: super_admin, admin, teacher, student, candidate")
			return
		}
		user.Role = body.Role
	}

	if err := h.userRepo.Update(user); err != nil {
		utils.InternalError(c, "Failed to update user")
		return
	}
	utils.OK(c, "User updated", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}
	if err := h.userRepo.Delete(id); err != nil {
		utils.InternalError(c, "Failed to delete user")
		return
	}
	utils.OK(c, "User deleted", nil)
}

func isValidRole(role models.Role) bool {
	for _, r := range models.AllRoles {
		if role == r {
			return true
		}
	}
	return false
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
