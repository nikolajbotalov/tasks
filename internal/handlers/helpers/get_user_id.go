package helpers

import (
	"TaskFlow/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	ctxUserID := c.Param("id")

	if ctxUserID == "" {
		return uuid.Nil, domain.ErrGetUserID
	}

	userID, err := uuid.Parse(ctxUserID)
	if err != nil {
		return uuid.Nil, domain.ErrGetUserID
	}

	return userID, nil
}
