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

func UpdateUserRole(uc rolesUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := helpers.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": response.IncorrectUserID})
			return
		}

		var req updateUserRoleRequest

		if err = c.ShouldBindJSON(&req); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		if err = uc.UpdateUserRole(userID, req.Role); err != nil {
			if errors.Is(err, domain.ErrUnknownRole) {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": response.UnknownRole})
				return
			}

			if errors.Is(err, domain.ErrUserNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": response.UserNotFound})
				return
			}

			logger.Error("update user role", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": response.RoleUpdated})
	}
}
