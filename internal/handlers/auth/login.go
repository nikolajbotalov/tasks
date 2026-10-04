package auth

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/validation"
	authUC "TaskFlow/internal/use-cases/auth"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginUser(uc authUC.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest domain.LoginUserRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		tokens, err := uc.LoginUser(loginRequest.Email, loginRequest.Password)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": response.InvalidCredentials})
				fmt.Printf("incorrect user email %s by error: %v\n", loginRequest.Email, err)
				return
			}

			fmt.Printf("failed to login user with email %s by error: %v\n", loginRequest.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, domain.TokenPair{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
