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

// Logout godoc
// @Summary Logout current user
// @Description Revokes the refresh sessions and clears the refresh cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} response.MessageResponse "logged out"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /auth/logout [post]
func Logout(uc authUC.UseCases, logger *slog.Logger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, _ := c.Cookie(refreshCookieName)

		if refreshToken == "" {
			var req refreshRequest
			if err := c.ShouldBindJSON(&req); err == nil {
				refreshToken = req.RefreshToken
			}
		}

		if refreshToken != "" {
			if err := uc.Logout(c.Request.Context(), refreshToken); err != nil {
				if !errors.Is(err, domain.ErrInvalidRefreshToken) &&
					!errors.Is(err, domain.ErrUserNotFound) {
					logger.Error("logout", slog.String("error", err.Error()))
					c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: response.Internal})
					return
				}
			}
		}

		clearRefreshCookie(c, cfg)
		c.JSON(http.StatusOK, response.MessageResponse{Message: response.LoggedOut})
	}
}
