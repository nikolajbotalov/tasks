package helpers

import (
	"TaskFlow/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetTaskID(c *gin.Context) (uuid.UUID, error) {
	ctxID := c.Param("id")
	if ctxID == "" {
		return uuid.Nil, domain.ErrGetTaskID
	}

	id, err := uuid.Parse(ctxID)
	if err != nil {
		return uuid.Nil, domain.ErrGetTaskID
	}

	return id, nil
}
