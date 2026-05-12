package api

import (
	"time"

	"github.com/alvindashahrul/my-app/internal/model"
)

type CreateStudentRequest struct {
	UserID         string  `json:"user_id" validate:"required,uuid"`
	NIS            string  `json:"nis" validate:"required,min=3,max=50"`
	NISN           *string `json:"nisn,omitempty" validate:"omitempty,min=3,max=50"`
	Gender         *string `json:"gender,omitempty" validate:"omitempty,oneof=MALE FEMALE"`
	BirthPlace     *string `json:"birth_place,omitempty" validate:"omitempty,max=100"`
	BirthDate      *string `json:"birth_date,omitempty" validate:"omitempty"`
	Religion       *string `json:"religion,omitempty" validate:"omitempty,max=50"`
	PhoneNumber    *string `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Address        *string `json:"address,omitempty"`
	PreviousSchool *string `json:"previous_school,omitempty" validate:"omitempty,max=150"`
	FatherName     *string `json:"father_name,omitempty" validate:"omitempty,max=150"`
	MotherName     *string `json:"mother_name,omitempty" validate:"omitempty,max=150"`
	ParentPhone    *string `json:"parent_phone,omitempty" validate:"omitempty,max=20"`
	PhotoURL       *string `json:"photo_url,omitempty"`
}

type UpdateStudentRequest struct {
	NIS            *string `json:"nis,omitempty" validate:"omitempty,min=3,max=50"`
	NISN           *string `json:"nisn,omitempty" validate:"omitempty,min=3,max=50"`
	Gender         *string `json:"gender,omitempty" validate:"omitempty,oneof=MALE FEMALE"`
	BirthPlace     *string `json:"birth_place,omitempty" validate:"omitempty,max=100"`
	BirthDate      *string `json:"birth_date,omitempty" validate:"omitempty"`
	Religion       *string `json:"religion,omitempty" validate:"omitempty,max=50"`
	PhoneNumber    *string `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Address        *string `json:"address,omitempty"`
	PreviousSchool *string `json:"previous_school,omitempty" validate:"omitempty,max=150"`
	FatherName     *string `json:"father_name,omitempty" validate:"omitempty,max=150"`
	MotherName     *string `json:"mother_name,omitempty" validate:"omitempty,max=150"`
	ParentPhone    *string `json:"parent_phone,omitempty" validate:"omitempty,max=20"`
	PhotoURL       *string `json:"photo_url,omitempty"`
	Status         *string `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE INACTIVE GRADUATED"`
}

type StudentResponse struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	NIS            string     `json:"nis"`
	NISN           *string    `json:"nisn,omitempty"`
	Gender         *string    `json:"gender,omitempty"`
	BirthPlace     *string    `json:"birth_place,omitempty"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`
	Religion       *string    `json:"religion,omitempty"`
	PhoneNumber    *string    `json:"phone_number,omitempty"`
	Address        *string    `json:"address,omitempty"`
	PreviousSchool *string    `json:"previous_school,omitempty"`
	FatherName     *string    `json:"father_name,omitempty"`
	MotherName     *string    `json:"mother_name,omitempty"`
	ParentPhone    *string    `json:"parent_phone,omitempty"`
	PhotoURL       *string    `json:"photo_url,omitempty"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type StudentWithUserResponse struct {
	StudentResponse
	User model.User `json:"user"`
}
