package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	authUC "TaskFlow/internal/use-cases/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken godoc
// @Summary Update user tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body refreshRequest false "optional when refresh cookie is used"
// @Success 200 {object} domain.RefreshResponse
// @Failure 401 {object} response.ErrorResponse "invalid or expired refresh token"
// @Failure 409 {object} response.ErrorResponse "refresh already in progress"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /auth/refresh [post]
func RefreshToken(uc authUC.UseCases, logger *slog.Logger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, cookieErr := c.Cookie(refreshCookieName)
		if cookieErr != nil || refreshToken == "" {
			var req refreshRequest
			if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
				c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.InvalidRefreshToken})
				clearRefreshCookie(c, cfg)
				return
			}
			refreshToken = req.RefreshToken
		}

		tokens, err := uc.RefreshToken(c.Request.Context(), refreshToken)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRefreshInProgress):
				c.JSON(http.StatusConflict, response.ErrorResponse{Error: response.RefreshInProgress})
			case errors.Is(err, domain.ErrInvalidRefreshToken):
				clearRefreshCookie(c, cfg)
				c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.InvalidRefreshToken})
			default:
				logger.Error("refresh tokens", slog.String("error", err.Error()))
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: response.Internal})
			}
			return

		}

		setRefreshCookie(c, cfg, tokens.RefreshToken, tokens.RefreshExpiresAt)
		c.JSON(http.StatusOK, domain.RefreshResponse{
			AccessToken: tokens.AccessToken,
		})
	}
}
