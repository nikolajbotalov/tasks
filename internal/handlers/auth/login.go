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

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginUser godoc
// @Summary Login existing user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginUserRequest true "user request body"
// @Success 200 {object} domain.TokenPair
// @Failure 400 {object} response.ErrorResponse "invalid request body"
// @Failure 401 {object} response.ErrorResponse "invalid login or password"
// @Failure 422 {object} response.ValidationErrorResponse "validation error"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /auth/login [post]
func LoginUser(uc authUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest loginUserRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		tokens, err := uc.LoginUser(loginRequest.Email, loginRequest.Password)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.InvalidCredentials})
				return
			}

			logger.Error("login", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusOK, domain.TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
