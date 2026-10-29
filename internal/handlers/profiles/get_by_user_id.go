package profiles

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	profilesUC "TaskFlow/internal/use-cases/profiles"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfileByUserID godoc
// @Summary Get profile by user ID
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Param id path string true "user id" Format(uuid)
// @Success 200 {object} domain.GetProfileResponse
// @Failure 400 {object} response.ErrorResponse "invalid user id"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 404 {object} response.ErrorResponse "profile not found"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /profiles/{id} [get]
func GetProfileByUserID(uc profilesUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := helpers.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.IncorrectUserID})
			return
		}

		profile, err := uc.GetProfileByUserID(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				c.JSON(http.StatusNotFound, response.ErrorResponse{Error: response.UserNotFound})
				return
			}

			logger.Error("get profile", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: response.Internal})
			return
		}

		c.JSON(http.StatusOK, domain.GetProfileResponse{Profile: profile})
	}
}
