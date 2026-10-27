package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/pagination"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// GetAll godoc
// @Summary Get all user tasks
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param assignee_id query string false "filter by assignee id"
// @Param status query string false "task status" Enums(to_do, in_progress, in_review, done, on_hold, cancelled)
// @Param priority query string false "task priority" Enums(low, middle, high)
// @Param search query string false "search"
// @Param page query int false "page pagination param"
// @Param limit query int false "limit pagination param"
// @Success 200 {object} domain.GetAllTasksResponse
// @Failure 400 {object} response.ErrorResponse "invalid query parameter"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /tasks [get]
func GetAll(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID, err := helpers.GetCallerID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.Unauthorized})
			return
		}

		assigneeID, err := helpers.GetTaskAssigneeID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.InvalidTaskAssigneeIDParam})
			return
		}

		taskStatus, err := helpers.GetTaskStatus(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.InvalidTaskStatusParam})
			return
		}

		taskPriority, err := helpers.GetTaskPriority(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.InvalidTaskPriorityParam})
			return
		}

		search := c.Query("search")
		trimmedSearch := strings.TrimSpace(search)
		if utf8.RuneCountInString(trimmedSearch) > handlers.MaxSearchLength {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.SearchQueryIsTooLong})
			return
		}

		page, limit, err := pagination.GetPaginationParams(c)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidPageParam) {
				c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.InvalidPageParam})
				return
			}

			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.InvalidLimitParam})
			return
		}

		tasks, err := uc.GetAllTasks(authorID, assigneeID, taskStatus, taskPriority, trimmedSearch, page, limit)
		if err != nil {
			logger.Error("get task list", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
