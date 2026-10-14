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

func DeleteTask(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, authorID, err := helpers.GetTaskAndAuthorIDs(c)
		if err != nil {
			if errors.Is(err, domain.ErrGetTaskID) {
				c.JSON(http.StatusBadRequest, gin.H{"error": response.IncorrectTaskID})
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": response.IncorrectAuthorID})
			return
		}

		err = uc.DeleteTask(id, authorID)
		if err != nil {
			if errors.Is(err, domain.ErrTaskNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": response.TaskNotFound})
				return
			}

			logger.Error("delete task", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": response.TaskDeleted})
	}
}
