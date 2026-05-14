package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamStatus string

const (
	ExamStatusDraft     ExamStatus = "draft"
	ExamStatusScheduled ExamStatus = "scheduled"
	ExamStatusOngoing   ExamStatus = "ongoing"
	ExamStatusFinished  ExamStatus = "finished"
)

type Exam struct {
	ID               uuid.UUID      `json:"id"                gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title            string         `json:"title"             gorm:"not null"`
	SubjectID        uuid.UUID      `json:"subject_id"        gorm:"type:uuid;not null"`
	Subject          Subject        `json:"subject,omitempty" gorm:"foreignKey:SubjectID"`
	ClassID          uuid.UUID      `json:"class_id"          gorm:"type:uuid;not null"`
	Class            Class          `json:"class,omitempty"   gorm:"foreignKey:ClassID"`
	CreatedByID      uuid.UUID      `json:"created_by_id"     gorm:"type:uuid;not null"`
	CreatedBy        User           `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	StartTime        time.Time      `json:"start_time"        gorm:"not null"`
	EndTime          time.Time      `json:"end_time"          gorm:"not null"`
	DurationMinutes  int            `json:"duration_minutes"  gorm:"not null"`
	Status           ExamStatus     `json:"status"            gorm:"type:varchar(20);default:'draft'"`
	ShuffleQuestions bool           `json:"shuffle_questions" gorm:"default:false"`
	ShuffleOptions   bool           `json:"shuffle_options"   gorm:"default:false"`
	AntiCheat        bool           `json:"anti_cheat"        gorm:"default:false"`
	PassingScore     int            `json:"passing_score"     gorm:"default:60"`
	TotalQuestions   int            `json:"total_questions"   gorm:"default:0"`
	ExamQuestions    []ExamQuestion `json:"exam_questions,omitempty" gorm:"foreignKey:ExamID"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (e *Exam) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type ExamQuestion struct {
	ID         uuid.UUID `json:"id"          gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ExamID     uuid.UUID `json:"exam_id"     gorm:"type:uuid;not null"`
	QuestionID uuid.UUID `json:"question_id" gorm:"type:uuid;not null"`
	Question   Question  `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	OrderNum   int       `json:"order_num"   gorm:"default:0"`
}

func (eq *ExamQuestion) BeforeCreate(tx *gorm.DB) error {
	if eq.ID == uuid.Nil {
		eq.ID = uuid.New()
	}
	return nil
}

type Class struct {
	ID        uuid.UUID `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name"       gorm:"not null"`
	Grade     string    `json:"grade"      gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Class) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type StudentClass struct {
	ID        uuid.UUID `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StudentID uuid.UUID `json:"student_id" gorm:"type:uuid;not null"`
	Student   User      `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	ClassID   uuid.UUID `json:"class_id"   gorm:"type:uuid;not null"`
	Class     Class     `json:"class,omitempty"   gorm:"foreignKey:ClassID"`
	CreatedAt time.Time `json:"created_at"`
}

func (sc *StudentClass) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	return nil
}

type ExamRequest struct {
	Title            string      `json:"title"             binding:"required"`
	SubjectID        uuid.UUID   `json:"subject_id"        binding:"required"`
	ClassID          uuid.UUID   `json:"class_id"          binding:"required"`
	StartTime        time.Time   `json:"start_time"        binding:"required"`
	EndTime          time.Time   `json:"end_time"          binding:"required"`
	DurationMinutes  int         `json:"duration_minutes"  binding:"required"`
	ShuffleQuestions bool        `json:"shuffle_questions"`
	ShuffleOptions   bool        `json:"shuffle_options"`
	AntiCheat        bool        `json:"anti_cheat"`
	PassingScore     int         `json:"passing_score"`
	QuestionIDs      []uuid.UUID `json:"question_ids"      binding:"required"`
}
