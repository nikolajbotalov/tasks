package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetByID godoc
// @Summary Get task details
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "task id" Format(uuid)
// @Success 200 {object} domain.Task
// @Failure 400 {object} response.ErrorResponse "invalid task id"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 404 {object} response.ErrorResponse "task not found"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /tasks/{id} [get]
func GetByID(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
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

			logger.Error("get task", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}
