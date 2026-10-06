package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetByID(uc tasksUC.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetTaskID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": response.IncorrectTaskID})
			return
		}

		task, err := uc.GetTaskByID(id)
		if err != nil {
			if errors.Is(err, domain.ErrTaskNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": response.TaskNotFound})
				return
			}

			fmt.Printf("getting task %v, error: %v\n", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}
