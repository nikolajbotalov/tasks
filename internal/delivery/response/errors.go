package response

type ErrorResponse struct {
	Error string `json:"error"`
}

type FieldError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"must be a valid email"`
}

type ValidationErrorResponse struct {
	Error   string       `json:"error" example:"validation failed"`
	Details []FieldError `json:"details,omitempty"`
}
