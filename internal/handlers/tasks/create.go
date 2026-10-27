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
	AssigneeID  uuid.UUID `json:"assignee_id"`
	Status      string    `json:"status" binding:"omitempty,oneof=to_do in_progress in_review done on_hold cancelled"`
	Priority    string    `json:"priority" binding:"omitempty,oneof=low middle high"`
	Name        string    `json:"name" binding:"required,min=2,max=255"`
	Description string    `json:"description" binding:"max=10000"`
}

// CreateTask godoc
// @Summary Create a task
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createTaskRequest true "task request body"
// @Success 201 {object} domain.Task
// @Failure 400 {object} response.ErrorResponse "invalid request body"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 422 {object} response.ValidationErrorResponse "validation error"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /tasks [post]
func CreateTask(uc tasksUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID, err := helpers.GetCallerID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.Unauthorized})
			return
		}

		var req createTaskRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		task, err := uc.CreateTask(authorID, req.AssigneeID, req.Name, req.Description, req.Status, req.Priority)
		if err != nil {
			logger.Error("create task", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}
