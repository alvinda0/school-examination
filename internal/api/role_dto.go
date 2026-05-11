package api

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
