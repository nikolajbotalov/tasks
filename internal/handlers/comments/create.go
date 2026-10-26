package comments

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/validation"
	commentsUC "TaskFlow/internal/use-cases/comment"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCommentRequest struct {
	Text string `json:"text" binding:"required,min=2,max=10000"`
}

// CreateTaskComment godoc
// @Summary Create comment for task
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "task id" Format(uuid)
// @Param request body createCommentRequest true "comment request body"
// @Success 201 {object} domain.TaskComment
// @Failure 400 {object} response.ErrorResponse "invalid task id or request body"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 404 {object} response.ErrorResponse "task not found"
// @Failure 422 {object} response.ValidationErrorResponse "validation error"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /tasks/{id}/comments [post]
func CreateTaskComment(uc commentsUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, authorID, err := helpers.GetTaskAndAuthorIDs(c)
		if err != nil {
			if errors.Is(err, domain.ErrGetTaskID) {
				c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.IncorrectTaskID})
				return
			}

			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: response.IncorrectAuthorID})
			return
		}

		var req createCommentRequest
		if err = c.ShouldBindJSON(&req); err != nil {
			validation.HandleBindError(c, err)
			return
		}

		taskComment, err := uc.CreateComment(taskID, authorID, req.Text)
		if err != nil {
			if errors.Is(err, domain.ErrTaskNotFound) {
				c.JSON(http.StatusNotFound, response.ErrorResponse{Error: response.TaskNotFound})
				return
			}

			logger.Error("create task comment", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusCreated, taskComment)
	}
}
