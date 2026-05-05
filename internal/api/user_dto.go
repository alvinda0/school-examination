package api

type CreateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
