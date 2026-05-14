package repository

import (
	"database/sql"
	"time"

	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/google/uuid"
)

type SubjectRepository interface {
	GetAll() ([]model.Subject, error)
	GetAllWithTeachers() ([]model.SubjectWithTeachers, error)
	GetSubjectsByClassID(classID uuid.UUID) ([]model.SubjectWithTeachers, error)
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

func (r *subjectRepository) GetAllWithTeachers() ([]model.SubjectWithTeachers, error) {
	query := `
		SELECT 
			s.id, s.name, s.code, s.description, s.created_at, s.updated_at, s.deleted_at,
			t.id as teacher_id, u.full_name as teacher_name
		FROM subjects s
		LEFT JOIN teacher_subjects ts ON s.id = ts.subject_id
		LEFT JOIN teachers t ON ts.teacher_id = t.id AND t.deleted_at IS NULL
		LEFT JOIN users u ON t.user_id = u.id AND u.deleted_at IS NULL
		WHERE s.deleted_at IS NULL
		ORDER BY s.name ASC, u.full_name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjectsMap := make(map[string]*model.SubjectWithTeachers)
	var subjectOrder []string

	for rows.Next() {
		var subject model.Subject
		var teacherID sql.NullString
		var teacherName sql.NullString

		err := rows.Scan(
			&subject.ID,
			&subject.Name,
			&subject.Code,
			&subject.Description,
			&subject.CreatedAt,
			&subject.UpdatedAt,
			&subject.DeletedAt,
			&teacherID,
			&teacherName,
		)
		if err != nil {
			return nil, err
		}

		subjectIDStr := subject.ID.String()

		// Jika subject belum ada di map, tambahkan
		if _, exists := subjectsMap[subjectIDStr]; !exists {
			subjectsMap[subjectIDStr] = &model.SubjectWithTeachers{
				Subject:  subject,
				Teachers: []model.TeacherInfo{},
			}
			subjectOrder = append(subjectOrder, subjectIDStr)
		}

		// Tambahkan teacher jika ada
		if teacherID.Valid && teacherName.Valid {
			teacherUUID, err := uuid.Parse(teacherID.String)
			if err == nil {
				subjectsMap[subjectIDStr].Teachers = append(
					subjectsMap[subjectIDStr].Teachers,
					model.TeacherInfo{
						ID:   teacherUUID,
						Name: teacherName.String,
					},
				)
			}
		}
	}

	// Convert map ke slice dengan urutan yang benar
	var subjects []model.SubjectWithTeachers
	for _, subjectID := range subjectOrder {
		subjects = append(subjects, *subjectsMap[subjectID])
	}

	return subjects, nil
}

func (r *subjectRepository) GetSubjectsByClassID(classID uuid.UUID) ([]model.SubjectWithTeachers, error) {
	// Get subjects taught by the homeroom teacher of the given class
	query := `
		SELECT 
			s.id, s.name, s.code, s.description, s.created_at, s.updated_at, s.deleted_at,
			t.id as teacher_id, u.full_name as teacher_name
		FROM subjects s
		JOIN teacher_subjects ts ON s.id = ts.subject_id
		JOIN teachers t ON ts.teacher_id = t.id AND t.deleted_at IS NULL
		JOIN users u ON t.user_id = u.id AND u.deleted_at IS NULL
		JOIN classes c ON c.homeroom_teacher_id = t.id AND c.deleted_at IS NULL
		WHERE s.deleted_at IS NULL AND c.id = $1
		ORDER BY s.name ASC, u.full_name ASC
	`

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjectsMap := make(map[string]*model.SubjectWithTeachers)
	var subjectOrder []string

	for rows.Next() {
		var subject model.Subject
		var teacherID sql.NullString
		var teacherName sql.NullString

		err := rows.Scan(
			&subject.ID,
			&subject.Name,
			&subject.Code,
			&subject.Description,
			&subject.CreatedAt,
			&subject.UpdatedAt,
			&subject.DeletedAt,
			&teacherID,
			&teacherName,
		)
		if err != nil {
			return nil, err
		}

		subjectIDStr := subject.ID.String()

		if _, exists := subjectsMap[subjectIDStr]; !exists {
			subjectsMap[subjectIDStr] = &model.SubjectWithTeachers{
				Subject:  subject,
				Teachers: []model.TeacherInfo{},
			}
			subjectOrder = append(subjectOrder, subjectIDStr)
		}

		if teacherID.Valid && teacherName.Valid {
			teacherUUID, err := uuid.Parse(teacherID.String)
			if err == nil {
				subjectsMap[subjectIDStr].Teachers = append(
					subjectsMap[subjectIDStr].Teachers,
					model.TeacherInfo{
						ID:   teacherUUID,
						Name: teacherName.String,
					},
				)
			}
		}
	}

	var subjects []model.SubjectWithTeachers
	for _, subjectID := range subjectOrder {
		subjects = append(subjects, *subjectsMap[subjectID])
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
