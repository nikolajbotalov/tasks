package helpers

import (
	"TaskFlow/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetTaskAssigneeID(c *gin.Context) (*uuid.UUID, error) {
	ctxAssigneeID := c.Query(domain.TaskAssigneeIDParam)
	if ctxAssigneeID == "" {
		return nil, nil
	}

	assigneeID, err := uuid.Parse(ctxAssigneeID)
	if err != nil {
		return nil, domain.ErrInvalidAssigneeParam
	}

	return &assigneeID, nil
}
