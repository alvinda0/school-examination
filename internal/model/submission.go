package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionStatus string

const (
	SubmissionStatusInProgress SubmissionStatus = "in_progress"
	SubmissionStatusSubmitted  SubmissionStatus = "submitted"
	SubmissionStatusGraded     SubmissionStatus = "graded"
)

type ExamSubmission struct {
	ID              uuid.UUID        `json:"id"               gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ExamID          uuid.UUID        `json:"exam_id"          gorm:"type:uuid;not null"`
	Exam            Exam             `json:"exam,omitempty"   gorm:"foreignKey:ExamID"`
	StudentID       uuid.UUID        `json:"student_id"       gorm:"type:uuid;not null"`
	Student         User             `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	StartedAt       time.Time        `json:"started_at"`
	SubmittedAt     *time.Time       `json:"submitted_at,omitempty"`
	Status          SubmissionStatus `json:"status"           gorm:"type:varchar(20);default:'in_progress'"`
	TotalScore      float64          `json:"total_score"      gorm:"default:0"`
	MaxScore        float64          `json:"max_score"        gorm:"default:0"`
	Percentage      float64          `json:"percentage"       gorm:"default:0"`
	IsPassed        bool             `json:"is_passed"        gorm:"default:false"`
	IsAutoSubmitted bool             `json:"is_auto_submitted" gorm:"default:false"`
	Answers         []StudentAnswer  `json:"answers,omitempty" gorm:"foreignKey:SubmissionID"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func (es *ExamSubmission) BeforeCreate(tx *gorm.DB) error {
	if es.ID == uuid.Nil {
		es.ID = uuid.New()
	}
	return nil
}

type StudentAnswer struct {
	ID             uuid.UUID  `json:"id"              gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SubmissionID   uuid.UUID  `json:"submission_id"   gorm:"type:uuid;not null"`
	QuestionID     uuid.UUID  `json:"question_id"     gorm:"type:uuid;not null"`
	Question       Question   `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	SelectedOption *uuid.UUID `json:"selected_option,omitempty" gorm:"type:uuid"`
	EssayAnswer    string     `json:"essay_answer,omitempty"   gorm:"type:text"`
	IsCorrect      *bool      `json:"is_correct,omitempty"`
	Score          float64    `json:"score"           gorm:"default:0"`
	GradedByID     *uuid.UUID `json:"graded_by_id,omitempty" gorm:"type:uuid"`
}

func (sa *StudentAnswer) BeforeCreate(tx *gorm.DB) error {
	if sa.ID == uuid.Nil {
		sa.ID = uuid.New()
	}
	return nil
}

type AnswerRequest struct {
	QuestionID     uuid.UUID  `json:"question_id"     binding:"required"`
	SelectedOption *uuid.UUID `json:"selected_option"`
	EssayAnswer    string     `json:"essay_answer"`
}

type SubmitRequest struct {
	Answers []AnswerRequest `json:"answers" binding:"required"`
}

type GradeEssayRequest struct {
	AnswerID uuid.UUID `json:"answer_id" binding:"required"`
	Score    float64   `json:"score"     binding:"required"`
}

type ExamResult struct {
	SubmissionID uuid.UUID  `json:"submission_id"`
	StudentName  string     `json:"student_name"`
	StudentEmail string     `json:"student_email"`
	TotalScore   float64    `json:"total_score"`
	MaxScore     float64    `json:"max_score"`
	Percentage   float64    `json:"percentage"`
	IsPassed     bool       `json:"is_passed"`
	Status       string     `json:"status"`
	SubmittedAt  *time.Time `json:"submitted_at"`
}
