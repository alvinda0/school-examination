package repository

import (
	"database/sql"

	"github.com/alvindashahrul/my-app/internal/model"
)

type UserRepository interface {
	GetAll(roleID string) ([]model.User, error)
	GetByID(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByIDWithRole(id string) (*UserWithRole, error)
	Create(fullName, email, password, roleID string, status bool) (*model.User, error)
	Patch(id string, email *string, status *bool) (*model.User, error)
	Delete(id string) (int64, error)
	UpdateLastLogin(id string) error
}

// UserWithRole adalah struct untuk user dengan informasi role
type UserWithRole struct {
	UserID   string
	FullName string
	Email    string
	RoleName string
	RoleID   string
	Status   bool
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(roleID string) ([]model.User, error) {
	var rows *sql.Rows
	var err error

	query := `SELECT u.id, u.full_name, u.email, u.role_id, r.name as role_name, u.status, u.last_login, u.created_at, u.updated_at, u.deleted_at 
	          FROM users u
	          LEFT JOIN roles r ON u.role_id = r.id
	          WHERE u.deleted_at IS NULL`

	if roleID != "" {
		query += " AND u.role_id = $1"
		rows, err = r.db.Query(query, roleID)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var roleName sql.NullString
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.RoleID, &roleName, &u.Status, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
			return nil, err
		}
		if roleName.Valid {
			u.RoleName = roleName.String
		}
		users = append(users, u)
	}

	if users == nil {
		users = []model.User{}
	}

	return users, nil
}

func (r *userRepository) GetByID(id string) (*model.User, error) {
	var u model.User
	var roleName sql.NullString
	query := `SELECT u.id, u.full_name, u.email, u.password, u.role_id, r.name as role_name, u.status, u.last_login, u.created_at, u.updated_at, u.deleted_at 
	          FROM users u
	          LEFT JOIN roles r ON u.role_id = r.id
	          WHERE u.id = $1 AND u.deleted_at IS NULL`
	
	err := r.db.QueryRow(query, id).
		Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.RoleID, &roleName, &u.Status, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if roleName.Valid {
		u.RoleName = roleName.String
	}

	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var u model.User
	query := `SELECT id, full_name, email, password, role_id, status, last_login, created_at, updated_at, deleted_at 
	          FROM users 
	          WHERE email = $1 AND deleted_at IS NULL`
	
	err := r.db.QueryRow(query, email).
		Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.RoleID, &u.Status, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) Create(fullName, email, password, roleID string, status bool) (*model.User, error) {
	var newUser model.User
	var roleName sql.NullString
	
	query := `INSERT INTO users (full_name, email, password, role_id, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) 
	          RETURNING id, full_name, email, role_id, 
	          (SELECT name FROM roles WHERE id = $4) as role_name,
	          status, last_login, created_at, updated_at, deleted_at`
	
	err := r.db.QueryRow(query, fullName, email, password, roleID, status).
		Scan(&newUser.ID, &newUser.FullName, &newUser.Email, &newUser.RoleID, &roleName, &newUser.Status, 
			&newUser.LastLogin, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.DeletedAt)

	if err != nil {
		return nil, err
	}

	if roleName.Valid {
		newUser.RoleName = roleName.String
	}

	return &newUser, nil
}

func (r *userRepository) Patch(id string, email *string, status *bool) (*model.User, error) {
	// Ambil data user yang ada
	existingUser, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, nil
	}

	// Update hanya field yang diberikan
	if email != nil {
		existingUser.Email = *email
	}
	if status != nil {
		existingUser.Status = *status
	}

	// Simpan perubahan
	var updated model.User
	var roleName sql.NullString
	query := `UPDATE users 
	         SET email = $1, status = $2, updated_at = NOW() 
	         WHERE id = $3 AND deleted_at IS NULL 
	         RETURNING id, full_name, email, role_id, 
	         (SELECT name FROM roles WHERE id = role_id) as role_name,
	         status, last_login, created_at, updated_at, deleted_at`
	
	err = r.db.QueryRow(query, existingUser.Email, existingUser.Status, id).
		Scan(&updated.ID, &updated.FullName, &updated.Email, &updated.RoleID, &roleName, &updated.Status, 
			&updated.LastLogin, &updated.CreatedAt, &updated.UpdatedAt, &updated.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if roleName.Valid {
		updated.RoleName = roleName.String
	}

	return &updated, nil
}

func (r *userRepository) Delete(id string) (int64, error) {
	// Soft delete
	result, err := r.db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *userRepository) UpdateLastLogin(id string) error {
	_, err := r.db.Exec("UPDATE users SET last_login = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	return err
}

func (r *userRepository) GetByIDWithRole(id string) (*UserWithRole, error) {
	var userWithRole UserWithRole
	query := `SELECT u.id, u.full_name, u.email, r.name as role_name, u.role_id, u.status 
	          FROM users u
	          INNER JOIN roles r ON u.role_id = r.id
	          WHERE u.id = $1 AND u.deleted_at IS NULL`
	
	err := r.db.QueryRow(query, id).
		Scan(&userWithRole.UserID, &userWithRole.FullName, &userWithRole.Email, 
			&userWithRole.RoleName, &userWithRole.RoleID, &userWithRole.Status)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &userWithRole, nil
}
