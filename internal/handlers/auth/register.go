package auth

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/use-cases/auth"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUser(uc auth.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerUser domain.RegisterUserRequest
		if err := c.ShouldBindJSON(&registerUser); err != nil {
			handleBindError(c, err)
			return
		}

		newUUID := uuid.New()

		err := uc.CreateUser(newUUID.String(), registerUser.Username, registerUser.Email, registerUser.Password)
		if err != nil {
			if errors.Is(err, domain.ErrEmailExists) {
				fmt.Printf("invalidate data from user with email %s by reason: %v\n", registerUser.Email, err)
				c.JSON(http.StatusConflict, gin.H{"error": response.EmailExists})
				return
			}

			fmt.Printf("cannot create user with email %s by reason: %v\n", registerUser.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusCreated, domain.RegisterResponse{
			ID: newUUID.String(),
		})
	}
}
