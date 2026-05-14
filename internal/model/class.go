package model

import (
	"time"

	"github.com/google/uuid"
)

type Class struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	GradeLevel         int        `json:"grade_level" db:"grade_level"`
	AcademicYear       string     `json:"academic_year" db:"academic_year"`
	HomeroomTeacherID  *uuid.UUID `json:"homeroom_teacher_id,omitempty" db:"homeroom_teacher_id"`
	MaxStudents        int        `json:"max_students" db:"max_students"`
	Status             string     `json:"status" db:"status"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ClassWithTeacher struct {
	Class
	HomeroomTeacherName *string `json:"homeroom_teacher_name,omitempty"`
}

type ClassWithTeacherDetail struct {
	Class
	HomeroomTeacher *Teacher `json:"homeroom_teacher,omitempty"`
}

type ClassWithStudents struct {
	Class
	Students        []Student `json:"students"`
	CurrentStudents int       `json:"current_students"`
}

type StudentInClass struct {
	Student
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type ClassWithStudentsDetail struct {
	Class
	Students        []StudentInClass `json:"students"`
	CurrentStudents int              `json:"current_students"`
}
