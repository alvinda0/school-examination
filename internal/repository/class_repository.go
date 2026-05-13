package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/alvindashahrul/my-app/internal/model"
)

type ClassRepository interface {
	Create(ctx context.Context, class *model.Class) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Class, error)
	GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]model.Class, int, error)
	GetAllWithTeacher(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]model.ClassWithTeacher, int, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetWithTeacher(ctx context.Context, id uuid.UUID) (*model.ClassWithTeacherDetail, error)
	GetWithStudents(ctx context.Context, id uuid.UUID) (*model.ClassWithStudents, error)
	GetStudentCount(ctx context.Context, classID uuid.UUID) (int, error)
}

type classRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(ctx context.Context, class *model.Class) error {
	query := `
		INSERT INTO classes (name, grade_level, academic_year, homeroom_teacher_id, max_students, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		class.Name, class.GradeLevel, class.AcademicYear,
		class.HomeroomTeacherID, class.MaxStudents, class.Status,
	).Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt)
}

func (r *classRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Class, error) {
	var class model.Class
	query := `
		SELECT id, name, grade_level, academic_year, homeroom_teacher_id, 
		       max_students, status, created_at, updated_at, deleted_at
		FROM classes
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&class.ID, &class.Name, &class.GradeLevel, &class.AcademicYear,
		&class.HomeroomTeacherID, &class.MaxStudents, &class.Status,
		&class.CreatedAt, &class.UpdatedAt, &class.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &class, err
}

func (r *classRepository) GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]model.Class, int, error) {
	var classes []model.Class
	var total int

	whereClause := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argPos := 1

	if gradeLevel, ok := filters["grade_level"].(int); ok {
		whereClause = append(whereClause, fmt.Sprintf("grade_level = $%d", argPos))
		args = append(args, gradeLevel)
		argPos++
	}

	if academicYear, ok := filters["academic_year"].(string); ok {
		whereClause = append(whereClause, fmt.Sprintf("academic_year = $%d", argPos))
		args = append(args, academicYear)
		argPos++
	}

	if status, ok := filters["status"].(string); ok {
		whereClause = append(whereClause, fmt.Sprintf("status = $%d", argPos))
		args = append(args, status)
		argPos++
	}

	where := strings.Join(whereClause, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM classes WHERE %s", where)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, name, grade_level, academic_year, homeroom_teacher_id,
		       max_students, status, created_at, updated_at, deleted_at
		FROM classes
		WHERE %s
		ORDER BY grade_level ASC, name ASC
		LIMIT $%d OFFSET $%d
	`, where, argPos, argPos+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var class model.Class
		err := rows.Scan(
			&class.ID, &class.Name, &class.GradeLevel, &class.AcademicYear,
			&class.HomeroomTeacherID, &class.MaxStudents, &class.Status,
			&class.CreatedAt, &class.UpdatedAt, &class.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		classes = append(classes, class)
	}

	return classes, total, rows.Err()
}

func (r *classRepository) GetAllWithTeacher(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]model.ClassWithTeacher, int, error) {
	var classes []model.ClassWithTeacher
	var total int

	whereClause := []string{"c.deleted_at IS NULL"}
	args := []interface{}{}
	argPos := 1

	if gradeLevel, ok := filters["grade_level"].(int); ok {
		whereClause = append(whereClause, fmt.Sprintf("c.grade_level = $%d", argPos))
		args = append(args, gradeLevel)
		argPos++
	}

	if academicYear, ok := filters["academic_year"].(string); ok {
		whereClause = append(whereClause, fmt.Sprintf("c.academic_year = $%d", argPos))
		args = append(args, academicYear)
		argPos++
	}

	if status, ok := filters["status"].(string); ok {
		whereClause = append(whereClause, fmt.Sprintf("c.status = $%d", argPos))
		args = append(args, status)
		argPos++
	}

	where := strings.Join(whereClause, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM classes c WHERE %s", where)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results with teacher info
	query := fmt.Sprintf(`
		SELECT 
			c.id, c.name, c.grade_level, c.academic_year, c.homeroom_teacher_id,
			c.max_students, c.status, c.created_at, c.updated_at, c.deleted_at,
			t.id, t.user_id, t.nip, t.gender, t.birth_place, t.birth_date,
			t.religion, t.phone_number, t.address, t.photo_url, t.status,
			t.created_at, t.updated_at, t.deleted_at,
			u.id, u.full_name, u.email
		FROM classes c
		LEFT JOIN teachers t ON c.homeroom_teacher_id = t.id AND t.deleted_at IS NULL
		LEFT JOIN users u ON t.user_id = u.id AND u.deleted_at IS NULL
		WHERE %s
		ORDER BY c.academic_year DESC, c.grade_level ASC, c.name ASC
		LIMIT $%d OFFSET $%d
	`, where, argPos, argPos+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var classWithTeacher model.ClassWithTeacher
		var teacherID, teacherUserID sql.NullString
		var teacherNIP, teacherGender, teacherBirthPlace, teacherReligion sql.NullString
		var teacherPhone, teacherAddress, teacherPhotoURL, teacherStatus sql.NullString
		var teacherBirthDate sql.NullTime
		var teacherCreatedAt, teacherUpdatedAt, teacherDeletedAt sql.NullTime
		var userID, userFullName, userEmail sql.NullString

		err := rows.Scan(
			&classWithTeacher.ID, &classWithTeacher.Name, &classWithTeacher.GradeLevel,
			&classWithTeacher.AcademicYear, &classWithTeacher.HomeroomTeacherID,
			&classWithTeacher.MaxStudents, &classWithTeacher.Status,
			&classWithTeacher.CreatedAt, &classWithTeacher.UpdatedAt, &classWithTeacher.DeletedAt,
			&teacherID, &teacherUserID, &teacherNIP, &teacherGender,
			&teacherBirthPlace, &teacherBirthDate, &teacherReligion,
			&teacherPhone, &teacherAddress, &teacherPhotoURL, &teacherStatus,
			&teacherCreatedAt, &teacherUpdatedAt, &teacherDeletedAt,
			&userID, &userFullName, &userEmail,
		)
		if err != nil {
			return nil, 0, err
		}

		// Set teacher name only
		if userFullName.Valid {
			classWithTeacher.HomeroomTeacherName = &userFullName.String
		}

		classes = append(classes, classWithTeacher)
	}

	return classes, total, rows.Err()
}

func (r *classRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	updates["updated_at"] = time.Now()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argPos))
		args = append(args, value)
		argPos++
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE classes
		SET %s
		WHERE id = $%d AND deleted_at IS NULL
	`, strings.Join(setClauses, ", "), argPos)

	result, err := r.db.ExecContext(ctx, query, args...)
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

func (r *classRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE classes
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
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

func (r *classRepository) GetWithTeacher(ctx context.Context, id uuid.UUID) (*model.ClassWithTeacherDetail, error) {
	var result model.ClassWithTeacherDetail
	var teacher model.Teacher
	
	query := `
		SELECT 
			c.id, c.name, c.grade_level, c.academic_year, c.homeroom_teacher_id,
			c.max_students, c.status, c.created_at, c.updated_at, c.deleted_at,
			t.id, t.user_id, t.nip, t.gender, t.birth_place, t.birth_date,
			t.religion, t.phone_number, t.address, t.photo_url, t.status,
			t.created_at, t.updated_at, t.deleted_at
		FROM classes c
		LEFT JOIN teachers t ON c.homeroom_teacher_id = t.id AND t.deleted_at IS NULL
		WHERE c.id = $1 AND c.deleted_at IS NULL
	`
	
	var teacherID, teacherUserID sql.NullString
	var teacherNIP, teacherGender, teacherBirthPlace, teacherReligion sql.NullString
	var teacherPhone, teacherAddress, teacherPhotoURL, teacherStatus sql.NullString
	var teacherBirthDate sql.NullTime
	var teacherCreatedAt, teacherUpdatedAt sql.NullTime
	var teacherDeletedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.ID, &result.Name, &result.GradeLevel, &result.AcademicYear,
		&result.HomeroomTeacherID, &result.MaxStudents, &result.Status,
		&result.CreatedAt, &result.UpdatedAt, &result.DeletedAt,
		&teacherID, &teacherUserID, &teacherNIP, &teacherGender,
		&teacherBirthPlace, &teacherBirthDate, &teacherReligion,
		&teacherPhone, &teacherAddress, &teacherPhotoURL, &teacherStatus,
		&teacherCreatedAt, &teacherUpdatedAt, &teacherDeletedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// If teacher exists, populate it
	if teacherID.Valid {
		teacher.ID, _ = uuid.Parse(teacherID.String)
		teacher.UserID, _ = uuid.Parse(teacherUserID.String)
		teacher.NIP = &teacherNIP.String
		teacher.Gender = &teacherGender.String
		teacher.BirthPlace = &teacherBirthPlace.String
		if teacherBirthDate.Valid {
			teacher.BirthDate = &teacherBirthDate.Time
		}
		teacher.Religion = &teacherReligion.String
		teacher.PhoneNumber = &teacherPhone.String
		teacher.Address = &teacherAddress.String
		teacher.PhotoURL = &teacherPhotoURL.String
		teacher.Status = teacherStatus.String
		teacher.CreatedAt = teacherCreatedAt.Time
		teacher.UpdatedAt = teacherUpdatedAt.Time
		if teacherDeletedAt.Valid {
			teacher.DeletedAt = &teacherDeletedAt.Time
		}
		result.HomeroomTeacher = &teacher
	}
	
	return &result, nil
}

func (r *classRepository) GetWithStudents(ctx context.Context, id uuid.UUID) (*model.ClassWithStudents, error) {
	var result model.ClassWithStudents

	// Get class info
	classQuery := `
		SELECT id, name, grade_level, academic_year, homeroom_teacher_id,
		       max_students, status, created_at, updated_at, deleted_at
		FROM classes
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, classQuery, id).Scan(
		&result.ID, &result.Name, &result.GradeLevel, &result.AcademicYear,
		&result.HomeroomTeacherID, &result.MaxStudents, &result.Status,
		&result.CreatedAt, &result.UpdatedAt, &result.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get students
	studentsQuery := `
		SELECT id, user_id, class_id, nis, nisn, gender, birth_place, birth_date,
		       religion, phone_number, address, previous_school, father_name,
		       mother_name, parent_phone, photo_url, status, created_at, updated_at, deleted_at
		FROM students
		WHERE class_id = $1 AND deleted_at IS NULL
		ORDER BY nis ASC
	`
	rows, err := r.db.QueryContext(ctx, studentsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result.Students = []model.Student{}
	for rows.Next() {
		var student model.Student
		err := rows.Scan(
			&student.ID, &student.UserID, &student.ClassID, &student.NIS,
			&student.NISN, &student.Gender, &student.BirthPlace, &student.BirthDate,
			&student.Religion, &student.PhoneNumber, &student.Address,
			&student.PreviousSchool, &student.FatherName, &student.MotherName,
			&student.ParentPhone, &student.PhotoURL, &student.Status,
			&student.CreatedAt, &student.UpdatedAt, &student.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result.Students = append(result.Students, student)
	}

	result.CurrentStudents = len(result.Students)
	return &result, rows.Err()
}

func (r *classRepository) GetStudentCount(ctx context.Context, classID uuid.UUID) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM students
		WHERE class_id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, classID).Scan(&count)
	return count, err
}
