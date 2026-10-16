package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/validation"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateTaskRequest struct {
	AssigneeID  uuid.UUID `json:"assignee_id"`
	Status      *string   `json:"status,omitempty"`
	Priority    *string   `json:"priority,omitempty"`
	Name        *string   `json:"name,omitempty" binding:"min=2,max=255"`
	Description *string   `json:"description,omitempty" binding:"max=10000"`
}

func UpdateTask(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetTaskID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": response.IncorrectTaskID})
			return
		}

		var req UpdateTaskRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		err = uc.UpdateTask(id, req.AssigneeID, req.Name, req.Description, req.Status, req.Priority)
		if err != nil {
			if errors.Is(err, domain.ErrTaskNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": response.TaskNotFound})
				return
			}

			logger.Error("update task", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": response.TaskUpdated})
	}
}
