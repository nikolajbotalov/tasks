package roles

import (
	"TaskFlow/internal/delivery/response"
	rolesUC "TaskFlow/internal/use-cases/roles"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllUsersRoles godoc
// @Summary Get roles list
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.GetUserRolesResponse
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 403 {object} response.ErrorResponse "permission denied"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /users/roles [get]
func GetAllUsersRoles(uc rolesUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := uc.GetAllUsersRoles()
		if err != nil {
			logger.Error("get user roles", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusOK, roles)
	}
}
