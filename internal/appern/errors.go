package appern

import "errors"

var (
	ErrUserNotFound      = errors.New("user tidak ditemukan")
	ErrInvalidInput      = errors.New("input tidak valid")
	ErrEmptyName         = errors.New("name tidak boleh kosong")
	ErrInvalidAge        = errors.New("age harus lebih dari 0")
	ErrEmptyID           = errors.New("ID tidak boleh kosong")
	ErrDatabaseOperation = errors.New("database operation failed")
)
