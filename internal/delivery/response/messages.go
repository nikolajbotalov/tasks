package response

const (
	InvalidCredentials            = "invalid login or password"
	Internal                      = "internal server error"
	EmailExists                   = "email already exists"
	Unauthorized                  = "unauthorized"
	InvalidPageParam              = "invalid page param"
	InvalidLimitParam             = "invalid limit param"
	IncorrectTaskID               = "no task id or incorrect format"
	TaskNotFound                  = "task not found"
	TaskDeleted                   = "task was deleted"
	TaskUpdated                   = "task was updated"
	PermissionDenied              = "permission denied"
	InvalidTaskAssigneeIDParam    = "invalid assignee param"
	InvalidTaskStatusParam        = "invalid task status param"
	InvalidTaskPriorityParam      = "invalid task priority param"
	SearchQueryIsTooLong          = "search query is too long"
	InvalidRequestBody            = "invalid request body"
	AuthorizationHeaderMissing    = "authorization header missing"
	InvalidToken                  = "invalid token"
	InvalidTokenClaims            = "invalid token claims"
	RateLimitExceeded             = "rate limit exceeded"
	IncorrectUserID               = "no user id or incorrect format"
	UnknownRole                   = "unknown role"
	UserNotFound                  = "user not found"
	ProfileNotFound               = "profile not found"
	ProfileUpdated                = "profile was updated"
	ForeignProfileChangeForbidden = "changing foreign profile is not allowed"
	InvalidProfileData            = "invalid profile data"
	RoleUpdated                   = "role was updated"
	OwnRoleChangeForbidden        = "changing your own role is not allowed"
	LastAdminDemotingForbidden    = "cannot demote the last admin"
	InvalidRefreshToken           = "invalid refresh token"
	LoggedOut                     = "logged out"
)

type MessageResponse struct {
	Message string `json:"message"`
}
