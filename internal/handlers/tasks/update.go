package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/validation"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateTaskRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=255"`
	Description *string `json:"description" binding:"omitempty,max=10000"`
}

func UpdateTask(uc tasksUC.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetTaskID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": response.IncorrectTaskID})
			return
		}

		var taskRequest UpdateTaskRequest
		if err = c.ShouldBindJSON(&taskRequest); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		err = uc.UpdateTask(id, taskRequest.Name, taskRequest.Description)
		if err != nil {
			if errors.Is(err, domain.ErrTaskNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": response.TaskNotFound})
				return
			}

			fmt.Printf("update task %s: %v\n", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": response.TaskUpdated})
	}
}
