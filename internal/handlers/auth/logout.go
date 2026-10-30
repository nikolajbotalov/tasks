package auth

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	authUC "TaskFlow/internal/use-cases/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout godoc
// @Summary Logout current user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.MessageResponse "logged out"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 404 {object} response.ErrorResponse "user not found"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /auth/logout [post]
func Logout(uc authUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		callerID, err := helpers.GetCallerID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.Unauthorized})
			return
		}

		if err = uc.Logout(c.Request.Context(), callerID); err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				c.JSON(http.StatusNotFound, response.ErrorResponse{Error: response.UserNotFound})
				return
			}

			logger.Error("logout", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: response.Internal})
			return
		}

		c.JSON(http.StatusOK, response.MessageResponse{Message: response.LoggedOut})
	}
}
