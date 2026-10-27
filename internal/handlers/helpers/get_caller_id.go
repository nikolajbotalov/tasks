package helpers

import (
	"TaskFlow/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCallerID(c *gin.Context) (uuid.UUID, error) {
	ctxAuthorID, ok := c.Get(domain.JWTID)
	if !ok {
		return uuid.Nil, domain.ErrGetCallerID
	}

	id, ok := ctxAuthorID.(uuid.UUID)
	if !ok {
		return uuid.Nil, domain.ErrGetCallerID
	}

	return id, nil
}
