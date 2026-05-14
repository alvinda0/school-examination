package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/mapper"
	"github.com/alvindashahrul/my-app/internal/middleware"
	"github.com/alvindashahrul/my-app/internal/model"
	"github.com/alvindashahrul/my-app/internal/services"
	"github.com/alvindashahrul/my-app/internal/utils"
	"github.com/google/uuid"
)

type StudentHandler struct {
	studentService services.StudentService
}

func NewStudentHandler(studentService services.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// GET /api/v1/students
func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	students, total, err := h.studentService.GetAllStudents(r.Context(), page, pageSize)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response := mapper.ToStudentWithUserResponseList(students)
	meta := &api.Metadata{
		Page:       page,
		Limit:      pageSize,
		Total:      total,
		TotalPages: (total + pageSize - 1) / pageSize,
	}
	utils.JSONResponse(w, http.StatusOK, "Students retrieved successfully", response, meta)
}

// GET /api/v1/students/{id}
func (h *StudentHandler) GetStudentByID(w http.ResponseWriter, r *http.Request, id string) {
	studentID, err := uuid.Parse(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid student ID", nil, nil)
		return
	}

	student, err := h.studentService.GetStudentByID(r.Context(), studentID)
	if err != nil {
		utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
		return
	}

	response := mapper.ToStudentWithUserResponse(student)
	utils.JSONResponse(w, http.StatusOK, "Student retrieved successfully", response, nil)
}

// Route: /api/v1/students → GET, POST
func (h *StudentHandler) StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllStudents(w, r)
	case http.MethodPost:
		h.CreateStudent(w, r)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// Route: /api/v1/students/{id} → GET, PUT, DELETE
func (h *StudentHandler) StudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r, "/api/v1/students/")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "ID tidak boleh kosong", nil, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetStudentByID(w, r, id)
	case http.MethodPut:
		h.UpdateStudent(w, r, id)
	case http.MethodDelete:
		h.DeleteStudent(w, r, id)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// POST /api/v1/students
func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var req api.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body", nil, nil)
		return
	}

	// Parse user_id
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid user ID", nil, nil)
		return
	}

	// Create student model
	student := &model.Student{
		ID:             uuid.New(),
		UserID:         userID,
		NIS:            req.NIS,
		NISN:           req.NISN,
		Gender:         req.Gender,
		BirthPlace:     req.BirthPlace,
		Religion:       req.Religion,
		PhoneNumber:    req.PhoneNumber,
		Address:        req.Address,
		PreviousSchool: req.PreviousSchool,
		FatherName:     req.FatherName,
		MotherName:     req.MotherName,
		ParentPhone:    req.ParentPhone,
		PhotoURL:       req.PhotoURL,
		Status:         "ACTIVE",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Parse birth_date if provided
	if req.BirthDate != nil && *req.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, "Invalid birth_date format, use YYYY-MM-DD", nil, nil)
			return
		}
		student.BirthDate = &birthDate
	}

	// Create student
	if err := h.studentService.CreateStudent(r.Context(), student); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response := mapper.ToStudentResponse(student)

	entityID := student.ID
	r = middleware.SetAuditData(r, &middleware.AuditData{
		Action:     "create",
		EntityID:   &entityID,
		EntityType: "student",
		NewData: model.JSONB{
			"id":       student.ID,
			"user_id":  student.UserID,
			"nis":      student.NIS,
			"gender":   student.Gender,
			"status":   student.Status,
		},
	})

	utils.JSONResponse(w, http.StatusCreated, "Student created successfully", response, nil)
}

// PUT /api/v1/students/{id}
func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request, id string) {
	studentID, err := uuid.Parse(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid student ID", nil, nil)
		return
	}

	var req api.UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body", nil, nil)
		return
	}

	// Get existing student
	existingStudent, err := h.studentService.GetStudentByID(r.Context(), studentID)
	if err != nil {
		utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
		return
	}

	// Update fields if provided
	if req.NIS != nil {
		existingStudent.NIS = *req.NIS
	}
	if req.NISN != nil {
		existingStudent.NISN = req.NISN
	}
	if req.Gender != nil {
		existingStudent.Gender = req.Gender
	}
	if req.BirthPlace != nil {
		existingStudent.BirthPlace = req.BirthPlace
	}
	if req.BirthDate != nil && *req.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, "Invalid birth_date format, use YYYY-MM-DD", nil, nil)
			return
		}
		existingStudent.BirthDate = &birthDate
	}
	if req.Religion != nil {
		existingStudent.Religion = req.Religion
	}
	if req.PhoneNumber != nil {
		existingStudent.PhoneNumber = req.PhoneNumber
	}
	if req.Address != nil {
		existingStudent.Address = req.Address
	}
	if req.PreviousSchool != nil {
		existingStudent.PreviousSchool = req.PreviousSchool
	}
	if req.FatherName != nil {
		existingStudent.FatherName = req.FatherName
	}
	if req.MotherName != nil {
		existingStudent.MotherName = req.MotherName
	}
	if req.ParentPhone != nil {
		existingStudent.ParentPhone = req.ParentPhone
	}
	if req.PhotoURL != nil {
		existingStudent.PhotoURL = req.PhotoURL
	}
	if req.Status != nil {
		existingStudent.Status = *req.Status
	}

	existingStudent.UpdatedAt = time.Now()

	// Update student
	if err := h.studentService.UpdateStudent(r.Context(), studentID, &existingStudent.Student); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	response := mapper.ToStudentWithUserResponse(existingStudent)

	r = middleware.SetAuditData(r, &middleware.AuditData{
		Action:     "update",
		EntityID:   &studentID,
		EntityType: "student",
		Changes: model.JSONB{
			"nis":    existingStudent.NIS,
			"status": existingStudent.Status,
		},
	})

	utils.JSONResponse(w, http.StatusOK, "Student updated successfully", response, nil)
}

// DELETE /api/v1/students/{id}
func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request, id string) {
	studentID, err := uuid.Parse(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid student ID", nil, nil)
		return
	}

	// Ambil data sebelum dihapus
	existing, _ := h.studentService.GetStudentByID(r.Context(), studentID)

	if err := h.studentService.DeleteStudent(r.Context(), studentID); err != nil {
		utils.JSONResponse(w, http.StatusNotFound, err.Error(), nil, nil)
		return
	}

	auditData := &middleware.AuditData{
		Action:     "delete",
		EntityID:   &studentID,
		EntityType: "student",
	}
	if existing != nil {
		auditData.DeletedData = model.JSONB{
			"id":      existing.ID,
			"user_id": existing.UserID,
			"nis":     existing.NIS,
			"status":  existing.Status,
		}
	}
	r = middleware.SetAuditData(r, auditData)

	utils.JSONResponse(w, http.StatusOK, "Student deleted successfully", nil, nil)
}
