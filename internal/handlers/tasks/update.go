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

type updateTaskRequest struct {
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	Status      *string    `json:"status" binding:"omitempty,oneof=to_do in_progress in_review done on_hold cancelled"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low middle high"`
	Name        *string    `json:"name" binding:"omitempty,min=2,max=255"`
	Description *string    `json:"description" binding:"omitempty,max=10000"`
}

// UpdateTask godoc
// @Summary Update a task
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "task id" Format(uuid)
// @Param request body updateTaskRequest true "task request body"
// @Success 200 {object} response.MessageResponse "task was updated"
// @Failure 400 {object} response.ErrorResponse "invalid task id or request body"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 404 {object} response.ErrorResponse "task not found"
// @Failure 422 {object} response.ValidationErrorResponse "validation error"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /tasks/{id} [put]
func UpdateTask(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := helpers.GetTaskID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": response.IncorrectTaskID})
			return
		}

		var req updateTaskRequest
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
