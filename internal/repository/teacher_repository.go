package repository

import (
	"database/sql"
	"time"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/google/uuid"
)

type TeacherRepository interface {
	GetAll() ([]model.TeacherWithUser, error)
	GetByID(id string) (*model.TeacherWithSubjects, error)
	GetByUserID(userID string) (*model.Teacher, error)
	Create(teacher *model.Teacher) error
	Update(teacher *model.Teacher) error
	Delete(id string) error
	AssignSubjects(teacherID string, subjectIDs []string) error
	RemoveSubjects(teacherID string, subjectIDs []string) error
	GetTeacherSubjects(teacherID string) ([]model.Subject, error)
}

type teacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) TeacherRepository {
	return &teacherRepository{db: db}
}

func (r *teacherRepository) GetAll() ([]model.TeacherWithUser, error) {
	query := `
		SELECT 
			t.id, t.user_id, t.nip, t.gender, t.birth_place, t.birth_date,
			t.religion, t.phone_number, t.address, t.photo_url, t.status,
			t.created_at, t.updated_at, t.deleted_at,
			u.id, u.full_name, u.email, u.role_id, u.status, u.last_login,
			u.created_at, u.updated_at, u.deleted_at
		FROM teachers t
		INNER JOIN users u ON t.user_id = u.id
		WHERE t.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []model.TeacherWithUser
	for rows.Next() {
		var teacher model.TeacherWithUser
		err := rows.Scan(
			&teacher.ID,
			&teacher.UserID,
			&teacher.NIP,
			&teacher.Gender,
			&teacher.BirthPlace,
			&teacher.BirthDate,
			&teacher.Religion,
			&teacher.PhoneNumber,
			&teacher.Address,
			&teacher.PhotoURL,
			&teacher.Status,
			&teacher.CreatedAt,
			&teacher.UpdatedAt,
			&teacher.DeletedAt,
			&teacher.User.ID,
			&teacher.User.FullName,
			&teacher.User.Email,
			&teacher.User.RoleID,
			&teacher.User.Status,
			&teacher.User.LastLogin,
			&teacher.User.CreatedAt,
			&teacher.User.UpdatedAt,
			&teacher.User.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func (r *teacherRepository) GetByID(id string) (*model.TeacherWithSubjects, error) {
	query := `
		SELECT 
			t.id, t.user_id, t.nip, t.gender, t.birth_place, t.birth_date,
			t.religion, t.phone_number, t.address, t.photo_url, t.status,
			t.created_at, t.updated_at, t.deleted_at,
			u.id, u.full_name, u.email, u.role_id, u.status, u.last_login,
			u.created_at, u.updated_at, u.deleted_at
		FROM teachers t
		INNER JOIN users u ON t.user_id = u.id
		WHERE t.id = $1 AND t.deleted_at IS NULL AND u.deleted_at IS NULL
	`

	var teacher model.TeacherWithSubjects
	err := r.db.QueryRow(query, id).Scan(
		&teacher.ID,
		&teacher.UserID,
		&teacher.NIP,
		&teacher.Gender,
		&teacher.BirthPlace,
		&teacher.BirthDate,
		&teacher.Religion,
		&teacher.PhoneNumber,
		&teacher.Address,
		&teacher.PhotoURL,
		&teacher.Status,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
		&teacher.DeletedAt,
		&teacher.User.ID,
		&teacher.User.FullName,
		&teacher.User.Email,
		&teacher.User.RoleID,
		&teacher.User.Status,
		&teacher.User.LastLogin,
		&teacher.User.CreatedAt,
		&teacher.User.UpdatedAt,
		&teacher.User.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get subjects
	subjects, err := r.GetTeacherSubjects(id)
	if err != nil {
		return nil, err
	}
	teacher.Subjects = subjects

	return &teacher, nil
}

func (r *teacherRepository) GetByUserID(userID string) (*model.Teacher, error) {
	query := `
		SELECT id, user_id, nip, gender, birth_place, birth_date,
			   religion, phone_number, address, photo_url, status,
			   created_at, updated_at, deleted_at
		FROM teachers
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	var teacher model.Teacher
	err := r.db.QueryRow(query, userID).Scan(
		&teacher.ID,
		&teacher.UserID,
		&teacher.NIP,
		&teacher.Gender,
		&teacher.BirthPlace,
		&teacher.BirthDate,
		&teacher.Religion,
		&teacher.PhoneNumber,
		&teacher.Address,
		&teacher.PhotoURL,
		&teacher.Status,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
		&teacher.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &teacher, nil
}

func (r *teacherRepository) Create(teacher *model.Teacher) error {
	query := `
		INSERT INTO teachers (
			id, user_id, nip, gender, birth_place, birth_date,
			religion, phone_number, address, photo_url, status,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	teacher.ID = uuid.New()
	teacher.CreatedAt = time.Now()
	teacher.UpdatedAt = time.Now()

	if teacher.Status == "" {
		teacher.Status = "ACTIVE"
	}

	return r.db.QueryRow(
		query,
		teacher.ID,
		teacher.UserID,
		teacher.NIP,
		teacher.Gender,
		teacher.BirthPlace,
		teacher.BirthDate,
		teacher.Religion,
		teacher.PhoneNumber,
		teacher.Address,
		teacher.PhotoURL,
		teacher.Status,
		teacher.CreatedAt,
		teacher.UpdatedAt,
	).Scan(&teacher.ID, &teacher.CreatedAt, &teacher.UpdatedAt)
}

func (r *teacherRepository) Update(teacher *model.Teacher) error {
	query := `
		UPDATE teachers
		SET nip = $1, gender = $2, birth_place = $3, birth_date = $4,
			religion = $5, phone_number = $6, address = $7, photo_url = $8,
			status = $9, updated_at = $10
		WHERE id = $11 AND deleted_at IS NULL
		RETURNING updated_at
	`

	teacher.UpdatedAt = time.Now()

	return r.db.QueryRow(
		query,
		teacher.NIP,
		teacher.Gender,
		teacher.BirthPlace,
		teacher.BirthDate,
		teacher.Religion,
		teacher.PhoneNumber,
		teacher.Address,
		teacher.PhotoURL,
		teacher.Status,
		teacher.UpdatedAt,
		teacher.ID,
	).Scan(&teacher.UpdatedAt)
}

func (r *teacherRepository) Delete(id string) error {
	query := `
		UPDATE teachers
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

func (r *teacherRepository) AssignSubjects(teacherID string, subjectIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO teacher_subjects (id, teacher_id, subject_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (teacher_id, subject_id) DO NOTHING
	`

	for _, subjectID := range subjectIDs {
		_, err := tx.Exec(query, uuid.New(), teacherID, subjectID, time.Now())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *teacherRepository) RemoveSubjects(teacherID string, subjectIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		DELETE FROM teacher_subjects
		WHERE teacher_id = $1 AND subject_id = $2
	`

	for _, subjectID := range subjectIDs {
		_, err := tx.Exec(query, teacherID, subjectID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *teacherRepository) GetTeacherSubjects(teacherID string) ([]model.Subject, error) {
	query := `
		SELECT s.id, s.name, s.code, s.description, s.created_at, s.updated_at, s.deleted_at
		FROM subjects s
		INNER JOIN teacher_subjects ts ON s.id = ts.subject_id
		WHERE ts.teacher_id = $1 AND s.deleted_at IS NULL
		ORDER BY s.name ASC
	`

	rows, err := r.db.Query(query, teacherID)
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
