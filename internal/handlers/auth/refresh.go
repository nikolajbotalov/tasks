package auth

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/validation"
	authUC "TaskFlow/internal/use-cases/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken godoc
// @Summary Update user tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body refreshRequest true "request body"
// @Success 200 {object} domain.TokenPair
// @Failure 400 {object} response.ErrorResponse "invalid request body"
// @Failure 401 {object} response.ErrorResponse "invalid or expired refresh token"
// @Failure 422 {object} response.ValidationErrorResponse "validation error"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /auth/refresh [post]
func RefreshToken(uc authUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req refreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		tokens, err := uc.RefreshToken(c.Request.Context(), req.RefreshToken)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidRefreshToken) {
				c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.InvalidRefreshToken})
				return
			}

			logger.Error("refresh tokens", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: response.Internal})
			return
		}

		c.JSON(http.StatusOK, tokens)
	}
}
