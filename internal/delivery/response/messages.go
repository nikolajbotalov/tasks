package response

const (
	InvalidCredentials         = "invalid login or password"
	Internal                   = "internal server error"
	EmailExists                = "email already exists"
	Unauthorized               = "unauthorized"
	InvalidPageParam           = "invalid page param"
	InvalidLimitParam          = "invalid limit param"
	IncorrectTaskID            = "no task id or incorrect format"
	TaskNotFound               = "task not found"
	TaskDeleted                = "task was deleted"
	TaskUpdated                = "task was updated"
	PermissionDenied           = "permission denied"
	InvalidTaskAssigneeIDParam = "invalid assignee param"
	InvalidTaskStatusParam     = "invalid task status param"
	InvalidTaskPriorityParam   = "invalid task priority param"
	SearchQueryIsTooLong       = "search query is too long"
	InvalidRequestBody         = "invalid request body"
	AuthorizationHeaderMissing = "authorization header missing"
	InvalidToken               = "invalid token"
	InvalidTokenClaims         = "invalid token claims"
	InvalidTokenRole           = "invalid token role"
	RateLimitExceeded          = "rate limit exceeded"
	IncorrectUserID            = "no user id or incorrect format"
	UnknownRole                = "unknown role"
	UserNotFound               = "user not found"
	RoleUpdated                = "role was updated"
	OwnRoleChangeForbidden     = "changing your own role is not allowed"
	LastAdminDemotingForbidden = "cannot demote the last admin"
)

type MessageResponse struct {
	Message string `json:"message"`
}
