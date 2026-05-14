package handlers

import (
	"strconv"

	"school-examination/internal/middleware"
	"school-examination/internal/model"
	"school-examination/internal/repository"
	"school-examination/internal/services"
	"school-examination/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExamHandler struct {
	examService    *services.ExamService
	examRepo       *repository.ExamRepository
	submissionRepo *repository.SubmissionRepository
}

func NewExamHandler(examService *services.ExamService, examRepo *repository.ExamRepository, submissionRepo *repository.SubmissionRepository) *ExamHandler {
	return &ExamHandler{examService: examService, examRepo: examRepo, submissionRepo: submissionRepo}
}

func (h *ExamHandler) CreateExam(c *gin.Context) {
	var req model.ExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	userID := middleware.GetUserID(c)
	exam, err := h.examService.CreateExam(&req, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, "Exam created", exam)
}

func (h *ExamHandler) GetExams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var subjectID, classID uuid.UUID
	if s := c.Query("subject_id"); s != "" {
		subjectID, _ = uuid.Parse(s)
	}
	if s := c.Query("class_id"); s != "" {
		classID, _ = uuid.Parse(s)
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)

	exams, total, err := h.examRepo.FindAll(page, limit, subjectID, classID, userID, role)
	if err != nil {
		utils.InternalError(c, "Failed to fetch exams")
		return
	}
	utils.Paginated(c, "Exams fetched", exams, total, page, limit)
}

func (h *ExamHandler) GetExam(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid exam ID")
		return
	}
	exam, err := h.examRepo.FindByID(id)
	if err != nil {
		utils.NotFound(c, "Exam not found")
		return
	}
	utils.OK(c, "Exam fetched", exam)
}

func (h *ExamHandler) GetAvailableExams(c *gin.Context) {
	userID := middleware.GetUserID(c)
	h.examRepo.UpdateExpiredExams()
	exams, err := h.examRepo.FindAvailableForStudent(userID)
	if err != nil {
		utils.InternalError(c, "Failed to fetch available exams")
		return
	}
	utils.OK(c, "Available exams fetched", exams)
}

func (h *ExamHandler) StartExam(c *gin.Context) {
	examID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid exam ID")
		return
	}
	studentID := middleware.GetUserID(c)
	submission, questions, err := h.examService.StartExam(examID, studentID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, "Exam started", gin.H{
		"submission": submission,
		"questions":  questions,
	})
}

func (h *ExamHandler) SaveAnswer(c *gin.Context) {
	submissionID, err := uuid.Parse(c.Param("submission_id"))
	if err != nil {
		utils.BadRequest(c, "Invalid submission ID")
		return
	}
	var req model.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	studentID := middleware.GetUserID(c)
	if err := h.examService.SaveAnswer(submissionID, &req, studentID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, "Answer saved", nil)
}

func (h *ExamHandler) SubmitExam(c *gin.Context) {
	submissionID, err := uuid.Parse(c.Param("submission_id"))
	if err != nil {
		utils.BadRequest(c, "Invalid submission ID")
		return
	}
	var req model.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	studentID := middleware.GetUserID(c)
	submission, err := h.examService.SubmitExam(submissionID, studentID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, "Exam submitted successfully", submission)
}

func (h *ExamHandler) GetExamResults(c *gin.Context) {
	examID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid exam ID")
		return
	}
	results, err := h.examService.GetExamResults(examID)
	if err != nil {
		utils.InternalError(c, "Failed to fetch results")
		return
	}
	utils.OK(c, "Results fetched", results)
}

func (h *ExamHandler) GradeEssay(c *gin.Context) {
	var req model.GradeEssayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	gradedByID := middleware.GetUserID(c)
	if err := h.examService.GradeEssay(&req, gradedByID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, "Essay graded", nil)
}

func (h *ExamHandler) GetMyResults(c *gin.Context) {
	studentID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	submissions, total, err := h.submissionRepo.FindByStudent(studentID, page, limit)
	if err != nil {
		utils.InternalError(c, "Failed to fetch results")
		return
	}
	utils.Paginated(c, "My results fetched", submissions, total, page, limit)
}

// --- Class Handlers ---

func (h *ExamHandler) CreateClass(c *gin.Context) {
	var class model.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.examRepo.CreateClass(&class); err != nil {
		utils.InternalError(c, "Failed to create class")
		return
	}
	utils.Created(c, "Class created", class)
}

func (h *ExamHandler) GetClasses(c *gin.Context) {
	classes, err := h.examRepo.FindAllClasses()
	if err != nil {
		utils.InternalError(c, "Failed to fetch classes")
		return
	}
	utils.OK(c, "Classes fetched", classes)
}

func (h *ExamHandler) AssignStudentToClass(c *gin.Context) {
	var sc model.StudentClass
	if err := c.ShouldBindJSON(&sc); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.examRepo.AssignStudentToClass(&sc); err != nil {
		utils.InternalError(c, "Failed to assign student")
		return
	}
	utils.Created(c, "Student assigned to class", sc)
}

func (h *ExamHandler) GetStudentsByClass(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid class ID")
		return
	}
	students, err := h.examRepo.FindStudentsByClass(classID)
	if err != nil {
		utils.InternalError(c, "Failed to fetch students")
		return
	}
	utils.OK(c, "Students fetched", students)
}
