package repository

import (
	"database/sql"
	"time"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/google/uuid"
)

type SubjectRepository interface {
	GetAll() ([]model.Subject, error)
	GetByID(id string) (*model.Subject, error)
	Create(subject *model.Subject) error
	Update(subject *model.Subject) error
	Delete(id string) error
}

type subjectRepository struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) GetAll() ([]model.Subject, error) {
	query := `
		SELECT id, name, code, description, created_at, updated_at, deleted_at
		FROM subjects
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []model.Subject
	for rows.Next() {
		var subject model.Subject
		err := rows.Scan(
			&subject.ID,
			&subject.Name,
			&subject.Code,
			&subject.Description,
			&subject.CreatedAt,
			&subject.UpdatedAt,
			&subject.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (r *subjectRepository) GetByID(id string) (*model.Subject, error) {
	query := `
		SELECT id, name, code, description, created_at, updated_at, deleted_at
		FROM subjects
		WHERE id = $1 AND deleted_at IS NULL
	`

	var subject model.Subject
	err := r.db.QueryRow(query, id).Scan(
		&subject.ID,
		&subject.Name,
		&subject.Code,
		&subject.Description,
		&subject.CreatedAt,
		&subject.UpdatedAt,
		&subject.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &subject, nil
}

func (r *subjectRepository) Create(subject *model.Subject) error {
	query := `
		INSERT INTO subjects (id, name, code, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	subject.ID = uuid.New()
	subject.CreatedAt = time.Now()
	subject.UpdatedAt = time.Now()

	return r.db.QueryRow(
		query,
		subject.ID,
		subject.Name,
		subject.Code,
		subject.Description,
		subject.CreatedAt,
		subject.UpdatedAt,
	).Scan(&subject.ID, &subject.CreatedAt, &subject.UpdatedAt)
}

func (r *subjectRepository) Update(subject *model.Subject) error {
	query := `
		UPDATE subjects
		SET name = $1, code = $2, description = $3, updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING updated_at
	`

	subject.UpdatedAt = time.Now()

	return r.db.QueryRow(
		query,
		subject.Name,
		subject.Code,
		subject.Description,
		subject.UpdatedAt,
		subject.ID,
	).Scan(&subject.UpdatedAt)
}

func (r *subjectRepository) Delete(id string) error {
	query := `
		UPDATE subjects
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
