package roles

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/validation"
	rolesUC "TaskFlow/internal/use-cases/roles"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateUserRole godoc
// @Summary Update user role
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "user id" Format(uuid)
// @Param request body updateUserRoleRequest true "request body"
// @Success 200 {object} response.MessageResponse "role was updated"
// @Failure 400 {object} response.ErrorResponse "invalid user id or request body"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 403 {object} response.ErrorResponse "permission denied"
// @Failure 404 {object} response.ErrorResponse "user not found"
// @Failure 422 {object} response.ValidationErrorResponse "validation error or unknown role"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /users/{id}/roles [put]
func UpdateUserRole(uc rolesUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := helpers.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.IncorrectUserID})
			return
		}

		var req updateUserRoleRequest

		if err = c.ShouldBindJSON(&req); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		if err = uc.UpdateUserRole(userID, req.Role); err != nil {
			if errors.Is(err, domain.ErrUnknownRole) {
				c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{Error: response.UnknownRole})
				return
			}

			if errors.Is(err, domain.ErrUserNotFound) {
				c.JSON(http.StatusNotFound, response.ErrorResponse{Error: response.UserNotFound})
				return
			}

			logger.Error("update user role", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusOK, response.MessageResponse{Message: response.RoleUpdated})
	}
}
