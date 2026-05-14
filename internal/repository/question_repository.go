package repository

import (
	"school-examination/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(q *model.Question) error {
	return r.db.Create(q).Error
}

func (r *QuestionRepository) FindByID(id uuid.UUID) (*model.Question, error) {
	var q model.Question
	err := r.db.Preload("Options").Preload("Subject").Preload("CreatedBy").
		First(&q, "id = ?", id).Error
	return &q, err
}

func (r *QuestionRepository) FindAll(page, limit int, subjectID uuid.UUID, createdByID uuid.UUID, role model.Role) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{}).Preload("Options").Preload("Subject")

	if role == model.RoleTeacher && createdByID != uuid.Nil {
		query = query.Where("created_by_id = ?", createdByID)
	}
	if subjectID != uuid.Nil {
		query = query.Where("subject_id = ?", subjectID)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * limit).Limit(limit).Find(&questions).Error
	return questions, total, err
}

func (r *QuestionRepository) Update(q *model.Question) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(q).Error
}

func (r *QuestionRepository) Delete(id uuid.UUID) error {
	r.db.Where("question_id = ?", id).Delete(&model.Option{})
	return r.db.Delete(&model.Question{}, "id = ?", id).Error
}

// Subject CRUD
func (r *QuestionRepository) CreateSubject(s *model.Subject) error {
	return r.db.Create(s).Error
}

func (r *QuestionRepository) FindAllSubjects() ([]model.Subject, error) {
	var subjects []model.Subject
	err := r.db.Find(&subjects).Error
	return subjects, err
}

func (r *QuestionRepository) FindSubjectByID(id uuid.UUID) (*model.Subject, error) {
	var s model.Subject
	err := r.db.First(&s, "id = ?", id).Error
	return &s, err
}
