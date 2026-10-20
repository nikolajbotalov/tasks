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
				c.JSON(http.StatusConflict, gin.H{"error": response.EmailExists})
				return
			}

			logger.Error("create user", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusCreated, domain.RegisterResponse{
			ID: newUUID.String(),
		})
	}
}
