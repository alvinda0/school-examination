package repository

import (
	"database/sql"

	"github.com/alvindashahrul/my-app/internal/model"
)

type RoleRepository interface {
	GetAll() ([]model.Role, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetAll() ([]model.Role, error) {
	rows, err := r.db.Query("SELECT id, name FROM role")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if roles == nil {
		roles = []model.Role{}
	}

	return roles, nil
}
