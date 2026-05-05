package repository

import (
	"database/sql"

	"github.com/alvindashahrul/my-app/internal/model"
)

type UserRepository interface {
	GetAll(roleID string) ([]model.User, error)
	GetByID(id string) (*model.User, error)
	Create(name string, age int) (*model.User, error)
	Update(id, name string, age int) (*model.User, error)
	Delete(id string) (int64, error)
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

	if roleID != "" {
		rows, err = r.db.Query("SELECT id, name, age, role_id FROM users WHERE role_id = $1", roleID)
	} else {
		rows, err = r.db.Query("SELECT id, name, age, role_id FROM users")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.RoleID); err != nil {
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
	err := r.db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Age)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) Create(name string, age int) (*model.User, error) {
	var newUser model.User
	err := r.db.QueryRow(
		"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, name, age",
		name, age,
	).Scan(&newUser.ID, &newUser.Name, &newUser.Age)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *userRepository) Update(id, name string, age int) (*model.User, error) {
	var updated model.User
	err := r.db.QueryRow(
		"UPDATE users SET name = $1, age = $2 WHERE id = $3 RETURNING id, name, age",
		name, age, id,
	).Scan(&updated.ID, &updated.Name, &updated.Age)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *userRepository) Delete(id string) (int64, error) {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
