package api

type CreateSubjectRequest struct {
	Name        string  `json:"name" validate:"required"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateSubjectRequest struct {
	Name        *string `json:"name,omitempty"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
}
