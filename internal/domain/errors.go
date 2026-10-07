package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrInternal           = errors.New("internal server error")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrHashedPassword     = errors.New("hashed password error")
	ErrGetAuthorID        = errors.New("no author id or incorrect format")
	ErrGetTasks           = errors.New("tasks unavailable")
	ErrTaskNotFound       = errors.New("task not found")
	ErrGetTaskID          = errors.New("no task id or incorrect format")
	ErrGetComments        = errors.New("comments unavailable")
	ErrInvalidPageParam   = errors.New("invalid page param or type")
	ErrInvalidLimitParam  = errors.New("invalid limit param or type")
)
