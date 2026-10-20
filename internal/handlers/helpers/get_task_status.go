package helpers

import (
	"TaskFlow/internal/domain"

	"github.com/gin-gonic/gin"
)

func GetTaskStatus(c *gin.Context) (domain.TaskStatus, error) {
	ctxStatus := c.Query(domain.TaskStatusParam)
	if ctxStatus == "" {
		return "", nil
	}

	taskStatus := domain.TaskStatus(ctxStatus)
	if !taskStatus.IsValid() {
		return "", domain.ErrInvalidStatusParam
	}

	return taskStatus, nil
}
