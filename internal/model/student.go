package model

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	NIS            string     `json:"nis" db:"nis"`
	NISN           *string    `json:"nisn,omitempty" db:"nisn"`
	Gender         *string    `json:"gender,omitempty" db:"gender"`
	BirthPlace     *string    `json:"birth_place,omitempty" db:"birth_place"`
	BirthDate      *time.Time `json:"birth_date,omitempty" db:"birth_date"`
	Religion       *string    `json:"religion,omitempty" db:"religion"`
	PhoneNumber    *string    `json:"phone_number,omitempty" db:"phone_number"`
	Address        *string    `json:"address,omitempty" db:"address"`
	PreviousSchool *string    `json:"previous_school,omitempty" db:"previous_school"`
	FatherName     *string    `json:"father_name,omitempty" db:"father_name"`
	MotherName     *string    `json:"mother_name,omitempty" db:"mother_name"`
	ParentPhone    *string    `json:"parent_phone,omitempty" db:"parent_phone"`
	PhotoURL       *string    `json:"photo_url,omitempty" db:"photo_url"`
	Status         string     `json:"status" db:"status"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type StudentWithUser struct {
	Student
	User User `json:"user"`
}
