package repository

import (
	"time"

	"school-examination/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamRepository struct {
	db *gorm.DB
}

func NewExamRepository(db *gorm.DB) *ExamRepository {
	return &ExamRepository{db: db}
}

func (r *ExamRepository) Create(exam *model.Exam) error {
	return r.db.Create(exam).Error
}

func (r *ExamRepository) FindByID(id uuid.UUID) (*model.Exam, error) {
	var exam model.Exam
	err := r.db.Preload("Subject").Preload("Class").Preload("CreatedBy").
		Preload("ExamQuestions.Question.Options").
		First(&exam, "id = ?", id).Error
	return &exam, err
}

func (r *ExamRepository) FindAll(page, limit int, subjectID, classID, createdByID uuid.UUID, role model.Role) ([]model.Exam, int64, error) {
	var exams []model.Exam
	var total int64

	query := r.db.Model(&model.Exam{}).Preload("Subject").Preload("Class").Preload("CreatedBy")

	if role == model.RoleTeacher && createdByID != uuid.Nil {
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

func (r *ExamRepository) FindAvailableForStudent(studentID uuid.UUID) ([]model.Exam, error) {
	var exams []model.Exam
	now := time.Now()

	var studentClasses []model.StudentClass
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

func (r *ExamRepository) Update(exam *model.Exam) error {
	return r.db.Save(exam).Error
}

func (r *ExamRepository) Delete(id uuid.UUID) error {
	r.db.Where("exam_id = ?", id).Delete(&model.ExamQuestion{})
	return r.db.Delete(&model.Exam{}, "id = ?", id).Error
}

func (r *ExamRepository) AddQuestions(examID uuid.UUID, questionIDs []uuid.UUID) error {
	for i, qID := range questionIDs {
		eq := model.ExamQuestion{
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
	r.db.Model(&model.Exam{}).
		Where("status = ? AND start_time <= ?", model.ExamStatusScheduled, now).
		Update("status", model.ExamStatusOngoing)
	r.db.Model(&model.Exam{}).
		Where("status IN ? AND end_time <= ?", []string{"scheduled", "ongoing"}, now).
		Update("status", model.ExamStatusFinished)
	return nil
}

// Class CRUD
func (r *ExamRepository) CreateClass(c *model.Class) error {
	return r.db.Create(c).Error
}

func (r *ExamRepository) FindAllClasses() ([]model.Class, error) {
	var classes []model.Class
	err := r.db.Find(&classes).Error
	return classes, err
}

func (r *ExamRepository) FindClassByID(id uuid.UUID) (*model.Class, error) {
	var c model.Class
	err := r.db.First(&c, "id = ?", id).Error
	return &c, err
}

// StudentClass
func (r *ExamRepository) AssignStudentToClass(sc *model.StudentClass) error {
	return r.db.Create(sc).Error
}

func (r *ExamRepository) FindStudentsByClass(classID uuid.UUID) ([]model.User, error) {
	var students []model.User
	err := r.db.Joins("JOIN student_classes ON student_classes.student_id = users.id").
		Where("student_classes.class_id = ? AND users.role = ?", classID, model.RoleStudent).
		Find(&students).Error
	return students, err
}
