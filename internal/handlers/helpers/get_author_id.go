package helpers

import (
	"TaskFlow/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAuthorID(c *gin.Context) (uuid.UUID, error) {
	ctxAuthorID, ok := c.Get("id")
	if !ok {
		return uuid.Nil, domain.ErrGetAuthorID
	}

	id, ok := ctxAuthorID.(uuid.UUID)
	if !ok {
		return uuid.Nil, domain.ErrGetAuthorID
	}

	return id, nil
}
