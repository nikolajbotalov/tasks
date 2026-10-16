package helpers

import (
	"TaskFlow/internal/domain"

	"github.com/gin-gonic/gin"
)

func GetTaskPriority(c *gin.Context) (domain.TaskPriority, error) {
	ctxPriority := c.Query(domain.TaskPriorityParam)
	if ctxPriority == "" {
		return "", nil
	}

	taskPriority := domain.TaskPriority(ctxPriority)
	if !taskPriority.IsValid() {
		return "", domain.ErrInvalidPriorityParam
	}

	return taskPriority, nil
}
