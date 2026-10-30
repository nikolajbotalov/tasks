package comments

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/pagination"
	commentsUC "TaskFlow/internal/use-cases/comment"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllTaskComments godoc
// @Summary Get task comments
// @Tags Comments
// @Produce json
// @Security BearerAuth
// @Param id path string true "task id" Format(uuid)
// @Param page query int false "page pagination param"
// @Param limit query int false "limit pagination param"
// @Success 200 {object} domain.GetTaskCommentsResponse
// @Failure 400 {object} response.ErrorResponse "invalid task id, page or limit"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 429 {object} response.ErrorResponse "rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /tasks/{id}/comments [get]
func GetAllTaskComments(uc commentsUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := helpers.GetTaskID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: response.IncorrectTaskID})
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

		commentsRes, err := uc.GetCommentList(c.Request.Context(), taskID, page, limit)
		if err != nil {
			logger.Error("get comment list", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Internal,
			})
			return
		}

		c.JSON(http.StatusOK, commentsRes)
	}
}
