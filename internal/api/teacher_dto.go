package api

import "time"

type CreateTeacherRequest struct {
	UserID      string     `json:"user_id" validate:"required"`
	NIP         *string    `json:"nip,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	BirthPlace  *string    `json:"birth_place,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	Religion    *string    `json:"religion,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	Address     *string    `json:"address,omitempty"`
	PhotoURL    *string    `json:"photo_url,omitempty"`
	Status      *string    `json:"status,omitempty"`
}

type UpdateTeacherRequest struct {
	NIP         *string    `json:"nip,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	BirthPlace  *string    `json:"birth_place,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	Religion    *string    `json:"religion,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	Address     *string    `json:"address,omitempty"`
	PhotoURL    *string    `json:"photo_url,omitempty"`
	Status      *string    `json:"status,omitempty"`
}

type AssignSubjectRequest struct {
	SubjectIDs []string `json:"subject_ids" validate:"required"`
}
