package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/google/uuid"
)

type StudentRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Student, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Student, error)
	FindByIDWithUser(ctx context.Context, id uuid.UUID) (*model.StudentWithUser, error)
	FindAll(ctx context.Context, limit, offset int) ([]*model.Student, error)
	FindAllWithUser(ctx context.Context, limit, offset int) ([]*model.StudentWithUser, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, student *model.Student) error
	Update(ctx context.Context, student *model.Student) error
	UpdateFields(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
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
		SELECT id, user_id, class_id, nis, nisn, gender, birth_place, birth_date,
			   religion, phone_number, address, previous_school,
			   father_name, mother_name, parent_phone, photo_url, status,
			   created_at, updated_at, deleted_at
		FROM students
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID, &student.UserID, &student.ClassID, &student.NIS, &student.NISN, &student.Gender,
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

func (r *studentRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Student, error) {
	var student model.Student
	query := `
		SELECT id, user_id, class_id, nis, nisn, gender, birth_place, birth_date,
			   religion, phone_number, address, previous_school,
			   father_name, mother_name, parent_phone, photo_url, status,
			   created_at, updated_at, deleted_at
		FROM students
		WHERE user_id = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&student.ID, &student.UserID, &student.ClassID, &student.NIS, &student.NISN, &student.Gender,
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
		SELECT id, user_id, class_id, nis, nisn, gender, birth_place, birth_date,
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
			&student.ID, &student.UserID, &student.ClassID, &student.NIS, &student.NISN, &student.Gender,
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

func (r *studentRepository) FindByIDWithUser(ctx context.Context, id uuid.UUID) (*model.StudentWithUser, error) {
	var s model.StudentWithUser
	query := `
		SELECT s.id, s.user_id, s.class_id, s.nis, s.nisn, s.gender, s.birth_place, s.birth_date,
		       s.religion, s.phone_number, s.address, s.previous_school,
		       s.father_name, s.mother_name, s.parent_phone, s.photo_url, s.status,
		       s.created_at, s.updated_at, s.deleted_at,
		       u.full_name, u.email
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.id = $1 AND s.deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.ClassID, &s.NIS, &s.NISN, &s.Gender,
		&s.BirthPlace, &s.BirthDate, &s.Religion, &s.PhoneNumber,
		&s.Address, &s.PreviousSchool, &s.FatherName, &s.MotherName,
		&s.ParentPhone, &s.PhotoURL, &s.Status, &s.CreatedAt,
		&s.UpdatedAt, &s.DeletedAt,
		&s.User.FullName, &s.User.Email,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("student not found")
	}
	return &s, err
}

func (r *studentRepository) FindAllWithUser(ctx context.Context, limit, offset int) ([]*model.StudentWithUser, error) {
	query := `
		SELECT s.id, s.user_id, s.class_id, s.nis, s.nisn, s.gender, s.birth_place, s.birth_date,
		       s.religion, s.phone_number, s.address, s.previous_school,
		       s.father_name, s.mother_name, s.parent_phone, s.photo_url, s.status,
		       s.created_at, s.updated_at, s.deleted_at,
		       u.full_name, u.email
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.deleted_at IS NULL
		ORDER BY s.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*model.StudentWithUser
	for rows.Next() {
		var s model.StudentWithUser
		err := rows.Scan(
			&s.ID, &s.UserID, &s.ClassID, &s.NIS, &s.NISN, &s.Gender,
			&s.BirthPlace, &s.BirthDate, &s.Religion, &s.PhoneNumber,
			&s.Address, &s.PreviousSchool, &s.FatherName, &s.MotherName,
			&s.ParentPhone, &s.PhotoURL, &s.Status, &s.CreatedAt,
			&s.UpdatedAt, &s.DeletedAt,
			&s.User.FullName, &s.User.Email,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, &s)
	}
	return students, rows.Err()
}

func (r *studentRepository) Count(ctx context.Context) (int, error) {	var count int
	query := `SELECT COUNT(*) FROM students WHERE deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *studentRepository) Create(ctx context.Context, student *model.Student) error {
	query := `
		INSERT INTO students (
			id, user_id, nis, nisn, gender, birth_place, birth_date,
			religion, phone_number, address, previous_school,
			father_name, mother_name, parent_phone, photo_url, status,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		student.ID, student.UserID, student.NIS, student.NISN, student.Gender,
		student.BirthPlace, student.BirthDate, student.Religion, student.PhoneNumber,
		student.Address, student.PreviousSchool, student.FatherName, student.MotherName,
		student.ParentPhone, student.PhotoURL, student.Status, student.CreatedAt,
		student.UpdatedAt,
	)
	return err
}

func (r *studentRepository) Update(ctx context.Context, student *model.Student) error {
	query := `
		UPDATE students SET
			nis = $1, nisn = $2, gender = $3, birth_place = $4, birth_date = $5,
			religion = $6, phone_number = $7, address = $8, previous_school = $9,
			father_name = $10, mother_name = $11, parent_phone = $12, photo_url = $13,
			status = $14, updated_at = $15
		WHERE id = $16 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		student.NIS, student.NISN, student.Gender, student.BirthPlace, student.BirthDate,
		student.Religion, student.PhoneNumber, student.Address, student.PreviousSchool,
		student.FatherName, student.MotherName, student.ParentPhone, student.PhotoURL,
		student.Status, student.UpdatedAt, student.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

func (r *studentRepository) UpdateFields(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	// Build dynamic UPDATE query
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argPos))
		args = append(args, value)
		argPos++
	}

	// Always update updated_at
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now())
	argPos++

	// Add WHERE clause
	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE students SET %s
		WHERE id = $%d AND deleted_at IS NULL
	`, fmt.Sprintf("%s", setClauses[0]), argPos)

	// Build full query with all SET clauses
	for i := 1; i < len(setClauses); i++ {
		query = fmt.Sprintf(`
			UPDATE students SET %s
			WHERE id = $%d AND deleted_at IS NULL
		`, fmt.Sprintf("%s, %s", setClauses[0], setClauses[i]), argPos)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

func (r *studentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE students SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}
