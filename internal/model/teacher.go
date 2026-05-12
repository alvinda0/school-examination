package model

import (
	"time"

	"github.com/google/uuid"
)

type Teacher struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	NIP         *string    `json:"nip,omitempty" db:"nip"`
	Gender      *string    `json:"gender,omitempty" db:"gender"`
	BirthPlace  *string    `json:"birth_place,omitempty" db:"birth_place"`
	BirthDate   *time.Time `json:"birth_date,omitempty" db:"birth_date"`
	Religion    *string    `json:"religion,omitempty" db:"religion"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number"`
	Address     *string    `json:"address,omitempty" db:"address"`
	PhotoURL    *string    `json:"photo_url,omitempty" db:"photo_url"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type TeacherWithUser struct {
	Teacher
	User User `json:"user"`
}

type TeacherWithSubjects struct {
	Teacher
	User     User      `json:"user"`
	Subjects []Subject `json:"subjects"`
}

type TeacherSubject struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TeacherID uuid.UUID `json:"teacher_id" db:"teacher_id"`
	SubjectID uuid.UUID `json:"subject_id" db:"subject_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
