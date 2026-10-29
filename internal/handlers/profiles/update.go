package profiles

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/validation"
	profilesUC "TaskFlow/internal/use-cases/profiles"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateProfileRequest struct {
	Username  *string `json:"username" binding:"omitempty,min=2,max=100"`
	FirstName *string `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" binding:"omitempty,min=2,max=100"`
}

// UpdateProfile godoc
// @Summary Update profile
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path string true "user id" Format(uuid)
// @Param request body updateProfileRequest true "request body"
// @Success 200 {object} response.MessageResponse "profile was update"
// @Failure 400 {object} response.ErrorResponse "invalid user or format id"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 404 {object} response.ErrorResponse "profile not found"
// @Failure 422 {object} response.ValidationErrorResponse "validation errors"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /profiles/{id} [put]
func UpdateProfile(uc profilesUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		callerID, err := helpers.GetCallerID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.Unauthorized})
			return
		}

		userID, err := helpers.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.IncorrectUserID})
			return
		}

		var req updateProfileRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		err = uc.UpdateProfile(c.Request.Context(), req.Username, req.FirstName, req.LastName, userID, callerID)
		if err != nil {
			if errors.Is(err, domain.ErrForeignProfileChange) {
				c.JSON(http.StatusForbidden, response.ErrorResponse{Error: response.ForeignProfileChangeForbidden})
				return
			}

			if errors.Is(err, domain.ErrInvalidProfileData) {
				c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{Error: response.InvalidProfileData})
				return
			}

			if errors.Is(err, domain.ErrProfileNotFound) {
				c.JSON(http.StatusNotFound, response.ErrorResponse{Error: response.ProfileNotFound})
				return
			}

			logger.Error("update profile", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: response.Internal})
			return
		}

		c.JSON(http.StatusOK, response.MessageResponse{Message: response.ProfileUpdated})
	}
}
