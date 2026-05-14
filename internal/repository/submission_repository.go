package repository

import (
	"time"

	"school-examination/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Create(s *models.ExamSubmission) error {
	return r.db.Create(s).Error
}

func (r *SubmissionRepository) FindByID(id uuid.UUID) (*models.ExamSubmission, error) {
	var s models.ExamSubmission
	err := r.db.Preload("Exam.Subject").Preload("Student").
		Preload("Answers.Question.Options").
		First(&s, "id = ?", id).Error
	return &s, err
}

func (r *SubmissionRepository) FindByExamAndStudent(examID, studentID uuid.UUID) (*models.ExamSubmission, error) {
	var s models.ExamSubmission
	err := r.db.Where("exam_id = ? AND student_id = ?", examID, studentID).
		Preload("Answers").
		First(&s).Error
	return &s, err
}

func (r *SubmissionRepository) FindByExam(examID uuid.UUID) ([]models.ExamSubmission, error) {
	var submissions []models.ExamSubmission
	err := r.db.Where("exam_id = ?", examID).
		Preload("Student").
		Find(&submissions).Error
	return submissions, err
}

func (r *SubmissionRepository) FindByStudent(studentID uuid.UUID, page, limit int) ([]models.ExamSubmission, int64, error) {
	var submissions []models.ExamSubmission
	var total int64

	r.db.Model(&models.ExamSubmission{}).Where("student_id = ?", studentID).Count(&total)
	err := r.db.Where("student_id = ?", studentID).
		Preload("Exam.Subject").
		Offset((page-1)*limit).Limit(limit).
		Find(&submissions).Error
	return submissions, total, err
}

func (r *SubmissionRepository) Update(s *models.ExamSubmission) error {
	return r.db.Save(s).Error
}

func (r *SubmissionRepository) SaveAnswer(answer *models.StudentAnswer) error {
	var existing models.StudentAnswer
	err := r.db.Where("submission_id = ? AND question_id = ?", answer.SubmissionID, answer.QuestionID).
		First(&existing).Error
	if err == nil {
		answer.ID = existing.ID
		return r.db.Save(answer).Error
	}
	return r.db.Create(answer).Error
}

func (r *SubmissionRepository) AutoSubmitExpired() error {
	now := time.Now()
	var submissions []models.ExamSubmission
	r.db.Joins("JOIN exams ON exams.id = exam_submissions.exam_id").
		Where("exam_submissions.status = ? AND exams.end_time <= ?", models.SubmissionStatusInProgress, now).
		Find(&submissions)
	for i := range submissions {
		submissions[i].Status = models.SubmissionStatusSubmitted
		submissions[i].IsAutoSubmitted = true
		t := now
		submissions[i].SubmittedAt = &t
		r.db.Save(&submissions[i])
	}
	return nil
}

func (r *SubmissionRepository) FindAnswerByID(id uuid.UUID) (*models.StudentAnswer, error) {
	var a models.StudentAnswer
	err := r.db.First(&a, "id = ?", id).Error
	return &a, err
}

func (r *SubmissionRepository) UpdateAnswer(a *models.StudentAnswer) error {
	return r.db.Save(a).Error
}
