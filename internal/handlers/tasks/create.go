package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/validation"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createTaskRequest struct {
	AssigneeID  uuid.UUID `json:"assignee_id,omitempty"`
	Status      string    `json:"status,omitempty"`
	Priority    string    `json:"priority,omitempty"`
	Name        string    `json:"name" binding:"required,min=2,max=255"`
	Description string    `json:"description" binding:"max=10000"`
}

func CreateTask(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID, err := helpers.GetAuthorID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": response.IncorrectAuthorID})
			return
		}

		var req createTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		task, err := uc.CreateTask(authorID, req.AssigneeID, req.Name, req.Description, req.Status, req.Priority)
		if err != nil {
			logger.Error("create task", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.CreateTask})
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}
