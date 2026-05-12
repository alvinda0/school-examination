package repository

import (
	"database/sql"

	"github.com/alvindashahrul/my-app/internal/model"
)

type RoleRepository interface {
	GetAll() ([]model.Role, error)
	GetByID(id string) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Create(name, description string) (*model.Role, error)
	Update(id, name, description string) (*model.Role, error)
	Delete(id string) (int64, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetAll() ([]model.Role, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at, updated_at FROM roles ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if roles == nil {
		roles = []model.Role{}
	}

	return roles, nil
}

func (r *roleRepository) GetByID(id string) (*model.Role, error) {
	var role model.Role
	query := "SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1"
	
	err := r.db.QueryRow(query, id).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	query := "SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1"
	
	err := r.db.QueryRow(query, name).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) Create(name, description string) (*model.Role, error) {
	var newRole model.Role
	query := `INSERT INTO roles (name, description, created_at, updated_at) 
	          VALUES ($1, $2, NOW(), NOW()) 
	          RETURNING id, name, description, created_at, updated_at`
	
	err := r.db.QueryRow(query, name, description).
		Scan(&newRole.ID, &newRole.Name, &newRole.Description, &newRole.CreatedAt, &newRole.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &newRole, nil
}

func (r *roleRepository) Update(id, name, description string) (*model.Role, error) {
	var updated model.Role
	query := `UPDATE roles 
	          SET name = $1, description = $2, updated_at = NOW() 
	          WHERE id = $3 
	          RETURNING id, name, description, created_at, updated_at`
	
	err := r.db.QueryRow(query, name, description, id).
		Scan(&updated.ID, &updated.Name, &updated.Description, &updated.CreatedAt, &updated.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *roleRepository) Delete(id string) (int64, error) {
	result, err := r.db.Exec("DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
