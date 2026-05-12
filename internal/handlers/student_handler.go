package handlers

import (
	"net/http"
	"strconv"

	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/mapper"
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

	response := mapper.ToStudentResponseList(students)
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

	response := mapper.ToStudentResponse(student)
	utils.JSONResponse(w, http.StatusOK, "Student retrieved successfully", response, nil)
}

// Route: /api/v1/students → GET
func (h *StudentHandler) StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllStudents(w, r)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}

// Route: /api/v1/students/{id} → GET
func (h *StudentHandler) StudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r, "/api/v1/students/")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "ID tidak boleh kosong", nil, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetStudentByID(w, r, id)
	default:
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil, nil)
	}
}
