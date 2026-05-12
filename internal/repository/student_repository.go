package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/google/uuid"
)

type StudentRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Student, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Student, error)
	Count(ctx context.Context) (int, error)
}

type studentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Student, error) {
	var student model.Student
	query := `
		SELECT id, user_id, nis, nisn, gender, birth_place, birth_date,
			   religion, phone_number, address, previous_school,
			   father_name, mother_name, parent_phone, photo_url, status,
			   created_at, updated_at, deleted_at
		FROM students
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID, &student.UserID, &student.NIS, &student.NISN, &student.Gender,
		&student.BirthPlace, &student.BirthDate, &student.Religion, &student.PhoneNumber,
		&student.Address, &student.PreviousSchool, &student.FatherName, &student.MotherName,
		&student.ParentPhone, &student.PhotoURL, &student.Status, &student.CreatedAt,
		&student.UpdatedAt, &student.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("student not found")
	}
	return &student, err
}

func (r *studentRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Student, error) {
	query := `
		SELECT id, user_id, nis, nisn, gender, birth_place, birth_date,
			   religion, phone_number, address, previous_school,
			   father_name, mother_name, parent_phone, photo_url, status,
			   created_at, updated_at, deleted_at
		FROM students
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*model.Student
	for rows.Next() {
		var student model.Student
		err := rows.Scan(
			&student.ID, &student.UserID, &student.NIS, &student.NISN, &student.Gender,
			&student.BirthPlace, &student.BirthDate, &student.Religion, &student.PhoneNumber,
			&student.Address, &student.PreviousSchool, &student.FatherName, &student.MotherName,
			&student.ParentPhone, &student.PhotoURL, &student.Status, &student.CreatedAt,
			&student.UpdatedAt, &student.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	return students, rows.Err()
}

func (r *studentRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM students WHERE deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
