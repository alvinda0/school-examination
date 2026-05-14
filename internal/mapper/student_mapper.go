package mapper

import (
	"github.com/alvindashahrul/my-app/internal/api"
	"github.com/alvindashahrul/my-app/internal/model"
)

func ToStudentResponse(student *model.Student) *api.StudentResponse {
	return &api.StudentResponse{
		ID:             student.ID.String(),
		UserID:         student.UserID.String(),
		NIS:            student.NIS,
		NISN:           student.NISN,
		Gender:         student.Gender,
		BirthPlace:     student.BirthPlace,
		BirthDate:      student.BirthDate,
		Religion:       student.Religion,
		PhoneNumber:    student.PhoneNumber,
		Address:        student.Address,
		PreviousSchool: student.PreviousSchool,
		FatherName:     student.FatherName,
		MotherName:     student.MotherName,
		ParentPhone:    student.ParentPhone,
		PhotoURL:       student.PhotoURL,
		Status:         student.Status,
		CreatedAt:      student.CreatedAt,
		UpdatedAt:      student.UpdatedAt,
	}
}

func ToStudentWithUserResponse(studentWithUser *model.StudentWithUser) *api.StudentWithUserResponse {
	return &api.StudentWithUserResponse{
		StudentResponse: *ToStudentResponse(&studentWithUser.Student),
		FullName:        studentWithUser.User.FullName,
		Email:           studentWithUser.User.Email,
	}
}

func ToStudentResponseList(students []*model.Student) []*api.StudentResponse {
	responses := make([]*api.StudentResponse, len(students))
	for i, student := range students {
		responses[i] = ToStudentResponse(student)
	}
	return responses
}

func ToStudentWithUserResponseList(students []*model.StudentWithUser) []*api.StudentWithUserResponse {
	responses := make([]*api.StudentWithUserResponse, len(students))
	for i, student := range students {
		responses[i] = ToStudentWithUserResponse(student)
	}
	return responses
}
