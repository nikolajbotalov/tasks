package domain

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid login or password")
	ErrEmailExists          = errors.New("email already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrProfileNotFound      = errors.New("profile not found")
	ErrForeignProfileChange = errors.New("foreign profile change not allowed")
	ErrInvalidProfileData   = errors.New("invalid profile data")
	ErrHashedPassword       = errors.New("hashed password")
	ErrGetCallerID          = errors.New("no caller id or incorrect format")
	ErrTaskNotFound         = errors.New("task not found")
	ErrGetTaskID            = errors.New("no task id or incorrect format")
	ErrInvalidPageParam     = errors.New("invalid page param or type")
	ErrInvalidLimitParam    = errors.New("invalid limit param or type")
	ErrUnknownRole          = errors.New("unknown role")
	ErrInvalidAssigneeParam = errors.New("invalid task assignee id param or type")
	ErrInvalidStatusParam   = errors.New("invalid task status param or type")
	ErrInvalidPriorityParam = errors.New("invalid task priority param or type")
	ErrGetUserID            = errors.New("no user id or incorrect format")
	ErrSelfRoleChange       = errors.New("own role change not allowed")
	ErrLastAdmin            = errors.New("demoting the last admin is not allowed")
)
