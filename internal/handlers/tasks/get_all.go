package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/pagination"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAll(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID, err := helpers.GetAuthorID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": response.IncorrectAuthorID})
			return
		}

		page, limit, err := pagination.GetPaginationParams(c)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidPageParam) {
				c.JSON(http.StatusBadRequest, gin.H{"error": response.InvalidPageParam})
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": response.InvalidLimitParam})
			return
		}

		tasks, err := uc.GetAllTasks(authorID, page, limit)
		if err != nil {
			logger.Error("get task list", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
