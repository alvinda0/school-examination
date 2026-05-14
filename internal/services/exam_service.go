package services

import (
	"errors"
	"math/rand"
	"time"

	"school-examination/internal/model"
	"school-examination/internal/repository"

	"github.com/google/uuid"
)

type ExamService struct {
	examRepo       *repository.ExamRepository
	questionRepo   *repository.QuestionRepository
	submissionRepo *repository.SubmissionRepository
}

func NewExamService(
	examRepo *repository.ExamRepository,
	questionRepo *repository.QuestionRepository,
	submissionRepo *repository.SubmissionRepository,
) *ExamService {
	return &ExamService{
		examRepo:       examRepo,
		questionRepo:   questionRepo,
		submissionRepo: submissionRepo,
	}
}

func (s *ExamService) CreateExam(req *model.ExamRequest, createdByID uuid.UUID) (*model.Exam, error) {
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	exam := &model.Exam{
		Title:            req.Title,
		SubjectID:        req.SubjectID,
		ClassID:          req.ClassID,
		CreatedByID:      createdByID,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		DurationMinutes:  req.DurationMinutes,
		Status:           model.ExamStatusScheduled,
		ShuffleQuestions: req.ShuffleQuestions,
		ShuffleOptions:   req.ShuffleOptions,
		AntiCheat:        req.AntiCheat,
		PassingScore:     req.PassingScore,
		TotalQuestions:   len(req.QuestionIDs),
	}

	if err := s.examRepo.Create(exam); err != nil {
		return nil, err
	}
	if err := s.examRepo.AddQuestions(exam.ID, req.QuestionIDs); err != nil {
		return nil, err
	}
	return s.examRepo.FindByID(exam.ID)
}

func (s *ExamService) StartExam(examID, studentID uuid.UUID) (*model.ExamSubmission, []model.ExamQuestion, error) {
	s.examRepo.UpdateExpiredExams()

	exam, err := s.examRepo.FindByID(examID)
	if err != nil {
		return nil, nil, errors.New("exam not found")
	}

	now := time.Now()
	if now.Before(exam.StartTime) {
		return nil, nil, errors.New("exam has not started yet")
	}
	if now.After(exam.EndTime) {
		return nil, nil, errors.New("exam has already ended")
	}

	// Kembalikan sesi yang sudah ada jika masih in_progress
	existing, err := s.submissionRepo.FindByExamAndStudent(examID, studentID)
	if err == nil && existing.ID != uuid.Nil {
		if existing.Status != model.SubmissionStatusInProgress {
			return nil, nil, errors.New("you have already submitted this exam")
		}
		return existing, exam.ExamQuestions, nil
	}

	submission := &model.ExamSubmission{
		ExamID:    examID,
		StudentID: studentID,
		StartedAt: now,
		Status:    model.SubmissionStatusInProgress,
		MaxScore:  calculateMaxScore(exam.ExamQuestions),
	}
	if err := s.submissionRepo.Create(submission); err != nil {
		return nil, nil, err
	}

	questions := exam.ExamQuestions
	if exam.ShuffleQuestions {
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}
	if exam.ShuffleOptions {
		for i := range questions {
			opts := questions[i].Question.Options
			rand.Shuffle(len(opts), func(a, b int) { opts[a], opts[b] = opts[b], opts[a] })
			questions[i].Question.Options = opts
		}
	}

	return submission, questions, nil
}

func (s *ExamService) SaveAnswer(submissionID uuid.UUID, req *model.AnswerRequest, studentID uuid.UUID) error {
	submission, err := s.submissionRepo.FindByID(submissionID)
	if err != nil {
		return errors.New("submission not found")
	}
	if submission.StudentID != studentID {
		return errors.New("forbidden")
	}
	if submission.Status != model.SubmissionStatusInProgress {
		return errors.New("exam already submitted")
	}

	answer := &model.StudentAnswer{
		SubmissionID:   submissionID,
		QuestionID:     req.QuestionID,
		SelectedOption: req.SelectedOption,
		EssayAnswer:    req.EssayAnswer,
	}
	return s.submissionRepo.SaveAnswer(answer)
}

func (s *ExamService) SubmitExam(submissionID, studentID uuid.UUID, req *model.SubmitRequest) (*model.ExamSubmission, error) {
	submission, err := s.submissionRepo.FindByID(submissionID)
	if err != nil {
		return nil, errors.New("submission not found")
	}
	if submission.StudentID != studentID {
		return nil, errors.New("forbidden")
	}
	if submission.Status != model.SubmissionStatusInProgress {
		return nil, errors.New("exam already submitted")
	}

	for _, ans := range req.Answers {
		answer := &model.StudentAnswer{
			SubmissionID:   submissionID,
			QuestionID:     ans.QuestionID,
			SelectedOption: ans.SelectedOption,
			EssayAnswer:    ans.EssayAnswer,
		}
		s.submissionRepo.SaveAnswer(answer)
	}

	totalScore := s.autoGrade(submissionID)

	now := time.Now()
	percentage := 0.0
	if submission.MaxScore > 0 {
		percentage = (totalScore / submission.MaxScore) * 100
	}

	exam, _ := s.examRepo.FindByID(submission.ExamID)
	submission.Status = model.SubmissionStatusSubmitted
	submission.SubmittedAt = &now
	submission.TotalScore = totalScore
	submission.Percentage = percentage
	submission.IsPassed = percentage >= float64(exam.PassingScore)

	if err := s.submissionRepo.Update(submission); err != nil {
		return nil, err
	}
	return submission, nil
}

func (s *ExamService) autoGrade(submissionID uuid.UUID) float64 {
	submission, err := s.submissionRepo.FindByID(submissionID)
	if err != nil {
		return 0
	}

	var totalScore float64
	for i := range submission.Answers {
		ans := &submission.Answers[i]
		q := ans.Question

		if q.Type == model.QuestionTypeEssay {
			continue
		}
		if ans.SelectedOption == nil {
			isCorrect := false
			ans.IsCorrect = &isCorrect
			s.submissionRepo.UpdateAnswer(ans)
			continue
		}

		isCorrect := false
		for _, opt := range q.Options {
			if opt.ID == *ans.SelectedOption && opt.IsCorrect {
				isCorrect = true
				break
			}
		}
		ans.IsCorrect = &isCorrect
		if isCorrect {
			ans.Score = float64(q.Points)
			totalScore += ans.Score
		}
		s.submissionRepo.UpdateAnswer(ans)
	}
	return totalScore
}

func (s *ExamService) GradeEssay(req *model.GradeEssayRequest, gradedByID uuid.UUID) error {
	answer, err := s.submissionRepo.FindAnswerByID(req.AnswerID)
	if err != nil {
		return errors.New("answer not found")
	}

	isCorrect := req.Score > 0
	answer.IsCorrect = &isCorrect
	answer.Score = req.Score
	answer.GradedByID = &gradedByID

	if err := s.submissionRepo.UpdateAnswer(answer); err != nil {
		return err
	}

	submission, err := s.submissionRepo.FindByID(answer.SubmissionID)
	if err != nil {
		return err
	}

	var totalScore float64
	for _, a := range submission.Answers {
		totalScore += a.Score
	}

	percentage := 0.0
	if submission.MaxScore > 0 {
		percentage = (totalScore / submission.MaxScore) * 100
	}

	exam, _ := s.examRepo.FindByID(submission.ExamID)
	submission.TotalScore = totalScore
	submission.Percentage = percentage
	submission.IsPassed = percentage >= float64(exam.PassingScore)
	submission.Status = model.SubmissionStatusGraded

	return s.submissionRepo.Update(submission)
}

func (s *ExamService) GetExamResults(examID uuid.UUID) ([]model.ExamResult, error) {
	submissions, err := s.submissionRepo.FindByExam(examID)
	if err != nil {
		return nil, err
	}

	results := make([]model.ExamResult, 0, len(submissions))
	for _, sub := range submissions {
		results = append(results, model.ExamResult{
			SubmissionID: sub.ID,
			StudentName:  sub.Student.Name,
			StudentEmail: sub.Student.Email,
			TotalScore:   sub.TotalScore,
			MaxScore:     sub.MaxScore,
			Percentage:   sub.Percentage,
			IsPassed:     sub.IsPassed,
			Status:       string(sub.Status),
			SubmittedAt:  sub.SubmittedAt,
		})
	}
	return results, nil
}

func calculateMaxScore(questions []model.ExamQuestion) float64 {
	var total float64
	for _, eq := range questions {
		total += float64(eq.Question.Points)
	}
	return total
}
