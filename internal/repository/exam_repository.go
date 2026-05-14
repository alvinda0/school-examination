package repository

import (
	"time"

	"school-examination/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamRepository struct {
	db *gorm.DB
}

func NewExamRepository(db *gorm.DB) *ExamRepository {
	return &ExamRepository{db: db}
}

func (r *ExamRepository) Create(exam *models.Exam) error {
	return r.db.Create(exam).Error
}

func (r *ExamRepository) FindByID(id uuid.UUID) (*models.Exam, error) {
	var exam models.Exam
	err := r.db.Preload("Subject").Preload("Class").Preload("CreatedBy").
		Preload("ExamQuestions.Question.Options").
		First(&exam, "id = ?", id).Error
	return &exam, err
}

func (r *ExamRepository) FindAll(page, limit int, subjectID, classID, createdByID uuid.UUID, role models.Role) ([]models.Exam, int64, error) {
	var exams []models.Exam
	var total int64

	query := r.db.Model(&models.Exam{}).Preload("Subject").Preload("Class").Preload("CreatedBy")

	if role == models.RoleTeacher && createdByID != uuid.Nil {
		query = query.Where("created_by_id = ?", createdByID)
	}
	if subjectID != uuid.Nil {
		query = query.Where("subject_id = ?", subjectID)
	}
	if classID != uuid.Nil {
		query = query.Where("class_id = ?", classID)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * limit).Limit(limit).Find(&exams).Error
	return exams, total, err
}

func (r *ExamRepository) FindAvailableForStudent(studentID uuid.UUID) ([]models.Exam, error) {
	var exams []models.Exam
	now := time.Now()

	var studentClasses []models.StudentClass
	r.db.Where("student_id = ?", studentID).Find(&studentClasses)

	classIDs := make([]uuid.UUID, 0, len(studentClasses))
	for _, sc := range studentClasses {
		classIDs = append(classIDs, sc.ClassID)
	}

	if len(classIDs) == 0 {
		return exams, nil
	}

	err := r.db.Preload("Subject").Preload("Class").
		Where("class_id IN ? AND start_time <= ? AND end_time >= ? AND status IN ?",
			classIDs, now, now, []string{"scheduled", "ongoing"}).
		Find(&exams).Error
	return exams, err
}

func (r *ExamRepository) Update(exam *models.Exam) error {
	return r.db.Save(exam).Error
}

func (r *ExamRepository) Delete(id uuid.UUID) error {
	r.db.Where("exam_id = ?", id).Delete(&models.ExamQuestion{})
	return r.db.Delete(&models.Exam{}, "id = ?", id).Error
}

func (r *ExamRepository) AddQuestions(examID uuid.UUID, questionIDs []uuid.UUID) error {
	for i, qID := range questionIDs {
		eq := models.ExamQuestion{
			ExamID:     examID,
			QuestionID: qID,
			OrderNum:   i + 1,
		}
		if err := r.db.Create(&eq).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ExamRepository) UpdateExpiredExams() error {
	now := time.Now()
	r.db.Model(&models.Exam{}).
		Where("status = ? AND start_time <= ?", models.ExamStatusScheduled, now).
		Update("status", models.ExamStatusOngoing)
	r.db.Model(&models.Exam{}).
		Where("status IN ? AND end_time <= ?", []string{"scheduled", "ongoing"}, now).
		Update("status", models.ExamStatusFinished)
	return nil
}

// Class CRUD
func (r *ExamRepository) CreateClass(c *models.Class) error {
	return r.db.Create(c).Error
}

func (r *ExamRepository) FindAllClasses() ([]models.Class, error) {
	var classes []models.Class
	err := r.db.Find(&classes).Error
	return classes, err
}

func (r *ExamRepository) FindClassByID(id uuid.UUID) (*models.Class, error) {
	var c models.Class
	err := r.db.First(&c, "id = ?", id).Error
	return &c, err
}

// StudentClass
func (r *ExamRepository) AssignStudentToClass(sc *models.StudentClass) error {
	return r.db.Create(sc).Error
}

func (r *ExamRepository) FindStudentsByClass(classID uuid.UUID) ([]models.User, error) {
	var students []models.User
	err := r.db.Joins("JOIN student_classes ON student_classes.student_id = users.id").
		Where("student_classes.class_id = ? AND users.role = ?", classID, models.RoleStudent).
		Find(&students).Error
	return students, err
}
