package auth

import (
	router "TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerUser domain.RegisterUserRequest
		if err := c.ShouldBindJSON(&registerUser); err != nil {
			handleBindError(c, err)
			return
		}

		newUUID := uuid.New()
		createdAt := time.Now()
		updatedAt := time.Now()
		trimEmail := strings.TrimSpace(registerUser.Email)
		lowerEmail := strings.ToLower(trimEmail)
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": router.Internal})
			fmt.Println(err)
			return
		}

		err = db.QueryRow("INSERT INTO users (id, email, username, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $6) ON CONFLICT (email) DO NOTHING RETURNING id",
			newUUID.String(), lowerEmail, registerUser.Username, string(hashedPass), createdAt, updatedAt).Scan(&newUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusConflict, gin.H{"error": router.EmailExists})
				fmt.Println(err)
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": router.Internal})
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusCreated, domain.RegisterResponse{
			ID: newUUID.String(),
		})
	}
}
