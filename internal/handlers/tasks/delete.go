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

// DeleteTask godoc
// @Summary Delete a task
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "task id" Format(uuid)
// @Success 200 {object} response.MessageResponse "task was deleted"
// @Failure 400 {object} response.ErrorResponse "invalid task id"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 404 {object} response.ErrorResponse "task not found"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /tasks/{id} [delete]
func DeleteTask(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, authorID, err := helpers.GetTaskAndAuthorIDs(c)
		if err != nil {
			if errors.Is(err, domain.ErrGetTaskID) {
				c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.IncorrectTaskID})
				return
			}

			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.IncorrectAuthorID})
			return
		}

		err = uc.DeleteTask(id, authorID)
		if err != nil {
			if errors.Is(err, domain.ErrTaskNotFound) {
				c.JSON(http.StatusNotFound, response.ErrorResponse{Error: response.TaskNotFound})
				return
			}

			logger.Error("delete task", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusOK, response.MessageResponse{
			Message: response.TaskDeleted,
		})
	}
}
