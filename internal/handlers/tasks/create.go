package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/validation"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTaskRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"max=10000"`
}

func CreateTask(uc tasksUC.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID, err := helpers.GetAuthorID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": response.IncorrectAuthorID})
			return
		}

		var requestTask createTaskRequest
		if err := c.ShouldBindJSON(&requestTask); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		task, err := uc.CreateTask(authorID, requestTask.Name, requestTask.Description)
		if err != nil {
			fmt.Printf("Create task %s, reason: %v\n", requestTask.Name, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.CreateTask})
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}
