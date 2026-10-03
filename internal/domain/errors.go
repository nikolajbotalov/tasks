package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrInternal           = errors.New("internal server error")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)
