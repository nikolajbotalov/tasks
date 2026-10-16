package helpers

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetTaskStatus(c *gin.Context) (domain.TaskStatus, error) {
	ctxStatus := c.Query(domain.TaskStatusParam)
	if ctxStatus == "" {
		return "", nil
	}

	fmt.Printf("ctxStatus: %s\n", ctxStatus)

	taskStatus := domain.TaskStatus(ctxStatus)
	fmt.Printf("taskStatus: %s\n", taskStatus)
	if !taskStatus.IsValid() {
		return "", domain.ErrInvalidStatusParam
	}

	return taskStatus, nil
}
