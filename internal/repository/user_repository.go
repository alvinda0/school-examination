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
	Update(id, fullName, email, password, roleID string, status bool) (*model.User, error)
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

	query := `SELECT id, full_name, email, role_id, status, last_login, created_at, updated_at, deleted_at 
	          FROM users 
	          WHERE deleted_at IS NULL`

	if roleID != "" {
		query += " AND role_id = $1"
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
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.RoleID, &u.Status, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
			return nil, err
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
	query := `SELECT id, full_name, email, password, role_id, status, last_login, created_at, updated_at, deleted_at 
	          FROM users 
	          WHERE id = $1 AND deleted_at IS NULL`
	
	err := r.db.QueryRow(query, id).
		Scan(&u.ID, &u.FullName, &u.Email, &u.Password, &u.RoleID, &u.Status, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
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
	
	query := `INSERT INTO users (full_name, email, password, role_id, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) 
	          RETURNING id, full_name, email, role_id, status, last_login, created_at, updated_at, deleted_at`
	
	err := r.db.QueryRow(query, fullName, email, password, roleID, status).
		Scan(&newUser.ID, &newUser.FullName, &newUser.Email, &newUser.RoleID, &newUser.Status, 
			&newUser.LastLogin, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.DeletedAt)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *userRepository) Update(id, fullName, email, password, roleID string, status bool) (*model.User, error) {
	var updated model.User
	var query string
	var err error

	if password != "" {
		query = `UPDATE users 
		         SET full_name = $1, email = $2, password = $3, role_id = $4, status = $5, updated_at = NOW() 
		         WHERE id = $6 AND deleted_at IS NULL 
		         RETURNING id, full_name, email, role_id, status, last_login, created_at, updated_at, deleted_at`
		err = r.db.QueryRow(query, fullName, email, password, roleID, status, id).
			Scan(&updated.ID, &updated.FullName, &updated.Email, &updated.RoleID, &updated.Status, 
				&updated.LastLogin, &updated.CreatedAt, &updated.UpdatedAt, &updated.DeletedAt)
	} else {
		query = `UPDATE users 
		         SET full_name = $1, email = $2, role_id = $3, status = $4, updated_at = NOW() 
		         WHERE id = $5 AND deleted_at IS NULL 
		         RETURNING id, full_name, email, role_id, status, last_login, created_at, updated_at, deleted_at`
		err = r.db.QueryRow(query, fullName, email, roleID, status, id).
			Scan(&updated.ID, &updated.FullName, &updated.Email, &updated.RoleID, &updated.Status, 
				&updated.LastLogin, &updated.CreatedAt, &updated.UpdatedAt, &updated.DeletedAt)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
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
