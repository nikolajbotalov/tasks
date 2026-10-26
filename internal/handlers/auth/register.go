package auth

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/validation"
	"TaskFlow/internal/use-cases/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type registerUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=2,max=32"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterUser godoc
// @Summary Create new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerUserRequest true "create new user request body"
// @Success 201 {object} domain.RegisterResponse
// @Failure 400 {object} response.ErrorResponse "invalid request body"
// @Failure 409 {object} response.ErrorResponse "email already exists"
// @Failure 422 {object} response.ValidationErrorResponse "validation error"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /auth/register [post]
func RegisterUser(uc auth.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerUser registerUserRequest
		if err := c.ShouldBindJSON(&registerUser); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		newUUID := uuid.New()

		err := uc.CreateUser(newUUID.String(), registerUser.Username, registerUser.Email, registerUser.Password)
		if err != nil {
			if errors.Is(err, domain.ErrEmailExists) {
				c.JSON(http.StatusConflict, response.ErrorResponse{
					Error: response.EmailExists,
				})
				return
			}

			logger.Error("create user", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusCreated, domain.RegisterResponse{
			ID: newUUID.String(),
		})
	}
}
