package repository

import (
	"school-examination/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(q *models.Question) error {
	return r.db.Create(q).Error
}

func (r *QuestionRepository) FindByID(id uuid.UUID) (*models.Question, error) {
	var q models.Question
	err := r.db.Preload("Options").Preload("Subject").Preload("CreatedBy").
		First(&q, "id = ?", id).Error
	return &q, err
}

func (r *QuestionRepository) FindAll(page, limit int, subjectID uuid.UUID, createdByID uuid.UUID, role models.Role) ([]models.Question, int64, error) {
	var questions []models.Question
	var total int64

	query := r.db.Model(&models.Question{}).Preload("Options").Preload("Subject")

	if role == models.RoleTeacher && createdByID != uuid.Nil {
		query = query.Where("created_by_id = ?", createdByID)
	}
	if subjectID != uuid.Nil {
		query = query.Where("subject_id = ?", subjectID)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * limit).Limit(limit).Find(&questions).Error
	return questions, total, err
}

func (r *QuestionRepository) Update(q *models.Question) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(q).Error
}

func (r *QuestionRepository) Delete(id uuid.UUID) error {
	r.db.Where("question_id = ?", id).Delete(&models.Option{})
	return r.db.Delete(&models.Question{}, "id = ?", id).Error
}

// Subject CRUD
func (r *QuestionRepository) CreateSubject(s *models.Subject) error {
	return r.db.Create(s).Error
}

func (r *QuestionRepository) FindAllSubjects() ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Find(&subjects).Error
	return subjects, err
}

func (r *QuestionRepository) FindSubjectByID(id uuid.UUID) (*models.Subject, error) {
	var s models.Subject
	err := r.db.First(&s, "id = ?", id).Error
	return &s, err
}
