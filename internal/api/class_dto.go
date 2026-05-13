package api

import "github.com/google/uuid"

type CreateClassRequest struct {
	Name              string     `json:"name" validate:"required,max=100"`
	GradeLevel        int        `json:"grade_level" validate:"required,min=1,max=12"`
	AcademicYear      string     `json:"academic_year" validate:"required,max=20"`
	HomeroomTeacherID *uuid.UUID `json:"homeroom_teacher_id,omitempty"`
	MaxStudents       int        `json:"max_students" validate:"required,min=1,max=100"`
}

type UpdateClassRequest struct {
	Name              *string    `json:"name,omitempty" validate:"omitempty,max=100"`
	GradeLevel        *int       `json:"grade_level,omitempty" validate:"omitempty,min=1,max=12"`
	AcademicYear      *string    `json:"academic_year,omitempty" validate:"omitempty,max=20"`
	HomeroomTeacherID *uuid.UUID `json:"homeroom_teacher_id,omitempty"`
	MaxStudents       *int       `json:"max_students,omitempty" validate:"omitempty,min=1,max=100"`
	Status            *string    `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE INACTIVE"`
}

type AssignStudentsToClassRequest struct {
	StudentIDs []uuid.UUID `json:"student_ids" validate:"required,min=1"`
}

type ClassQueryParams struct {
	GradeLevel   *int    `query:"grade_level"`
	AcademicYear *string `query:"academic_year"`
	Status       *string `query:"status"`
	Page         int     `query:"page"`
	Limit        int     `query:"limit"`
}
