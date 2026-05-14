package handlers

import (
	"strconv"

	"school-examination/internal/middleware"
	"school-examination/internal/models"
	"school-examination/internal/repository"
	"school-examination/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuestionHandler struct {
	questionRepo *repository.QuestionRepository
}

func NewQuestionHandler(questionRepo *repository.QuestionRepository) *QuestionHandler {
	return &QuestionHandler{questionRepo: questionRepo}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req models.QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if req.Points == 0 {
		req.Points = 1
	}

	userID := middleware.GetUserID(c)
	question := &models.Question{
		SubjectID:   req.SubjectID,
		CreatedByID: userID,
		Type:        req.Type,
		Content:     req.Content,
		ImageURL:    req.ImageURL,
		Points:      req.Points,
		Explanation: req.Explanation,
	}
	for _, opt := range req.Options {
		question.Options = append(question.Options, models.Option{
			Content:   opt.Content,
			IsCorrect: opt.IsCorrect,
		})
	}

	if err := h.questionRepo.Create(question); err != nil {
		utils.InternalError(c, "Failed to create question")
		return
	}
	utils.Created(c, "Question created", question)
}

func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var subjectID uuid.UUID
	if s := c.Query("subject_id"); s != "" {
		subjectID, _ = uuid.Parse(s)
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)

	questions, total, err := h.questionRepo.FindAll(page, limit, subjectID, userID, role)
	if err != nil {
		utils.InternalError(c, "Failed to fetch questions")
		return
	}
	utils.Paginated(c, "Questions fetched", questions, total, page, limit)
}

func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid question ID")
		return
	}
	question, err := h.questionRepo.FindByID(id)
	if err != nil {
		utils.NotFound(c, "Question not found")
		return
	}
	utils.OK(c, "Question fetched", question)
}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid question ID")
		return
	}
	question, err := h.questionRepo.FindByID(id)
	if err != nil {
		utils.NotFound(c, "Question not found")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)
	if role == models.RoleTeacher && question.CreatedByID != userID {
		utils.Forbidden(c, "You can only edit your own questions")
		return
	}

	var req models.QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	question.SubjectID = req.SubjectID
	question.Type = req.Type
	question.Content = req.Content
	question.ImageURL = req.ImageURL
	question.Points = req.Points
	question.Explanation = req.Explanation
	question.Options = nil
	for _, opt := range req.Options {
		question.Options = append(question.Options, models.Option{
			Content:   opt.Content,
			IsCorrect: opt.IsCorrect,
		})
	}

	if err := h.questionRepo.Update(question); err != nil {
		utils.InternalError(c, "Failed to update question")
		return
	}
	utils.OK(c, "Question updated", question)
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid question ID")
		return
	}
	question, err := h.questionRepo.FindByID(id)
	if err != nil {
		utils.NotFound(c, "Question not found")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)
	if role == models.RoleTeacher && question.CreatedByID != userID {
		utils.Forbidden(c, "You can only delete your own questions")
		return
	}

	if err := h.questionRepo.Delete(id); err != nil {
		utils.InternalError(c, "Failed to delete question")
		return
	}
	utils.OK(c, "Question deleted", nil)
}

func (h *QuestionHandler) CreateSubject(c *gin.Context) {
	var subject models.Subject
	if err := c.ShouldBindJSON(&subject); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := h.questionRepo.CreateSubject(&subject); err != nil {
		utils.InternalError(c, "Failed to create subject")
		return
	}
	utils.Created(c, "Subject created", subject)
}

func (h *QuestionHandler) GetSubjects(c *gin.Context) {
	subjects, err := h.questionRepo.FindAllSubjects()
	if err != nil {
		utils.InternalError(c, "Failed to fetch subjects")
		return
	}
	utils.OK(c, "Subjects fetched", subjects)
}
