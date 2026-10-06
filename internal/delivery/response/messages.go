package response

const (
	InvalidCredentials = "invalid login or password"
	Internal           = "internal server error"
	EmailExists        = "email already exists"
	IncorrectAuthorID  = "no author id or incorrect format"
	CreateTask         = "failed to create task"
	InvalidPageParam   = "invalid page param"
	InvalidLimitParam  = "invalid limit param"
)
