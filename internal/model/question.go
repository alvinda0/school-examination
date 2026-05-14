package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionType string

const (
	QuestionTypePG        QuestionType = "multiple_choice"
	QuestionTypeEssay     QuestionType = "essay"
	QuestionTypeTrueFalse QuestionType = "true_false"
)

type Question struct {
	ID          uuid.UUID    `json:"id"            gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SubjectID   uuid.UUID    `json:"subject_id"    gorm:"type:uuid;not null"`
	Subject     Subject      `json:"subject,omitempty"     gorm:"foreignKey:SubjectID"`
	CreatedByID uuid.UUID    `json:"created_by_id" gorm:"type:uuid;not null"`
	CreatedBy   User         `json:"created_by,omitempty"  gorm:"foreignKey:CreatedByID"`
	Type        QuestionType `json:"type"          gorm:"type:varchar(20);not null"`
	Content     string       `json:"content"       gorm:"type:text;not null"`
	ImageURL    string       `json:"image_url,omitempty"`
	Points      int          `json:"points"        gorm:"default:1"`
	Explanation string       `json:"explanation,omitempty" gorm:"type:text"`
	Options     []Option     `json:"options,omitempty"     gorm:"foreignKey:QuestionID"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

type Option struct {
	ID         uuid.UUID `json:"id"          gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	QuestionID uuid.UUID `json:"question_id" gorm:"type:uuid;not null"`
	Content    string    `json:"content"     gorm:"type:text;not null"`
	IsCorrect  bool      `json:"is_correct"  gorm:"default:false"`
}

func (o *Option) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type Subject struct {
	ID          uuid.UUID `json:"id"          gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name"        gorm:"not null"`
	Code        string    `json:"code"        gorm:"uniqueIndex;not null"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Subject) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type QuestionRequest struct {
	SubjectID   uuid.UUID       `json:"subject_id"  binding:"required"`
	Type        QuestionType    `json:"type"        binding:"required"`
	Content     string          `json:"content"     binding:"required"`
	ImageURL    string          `json:"image_url"`
	Points      int             `json:"points"`
	Explanation string          `json:"explanation"`
	Options     []OptionRequest `json:"options"`
}

type OptionRequest struct {
	Content   string `json:"content"    binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}
