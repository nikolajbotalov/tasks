package comments

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/pagination"
	commentsUC "TaskFlow/internal/use-cases/comment"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTaskComments(uc commentsUC.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := helpers.GetTaskID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": response.IncorrectTaskID})
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

		commentsRes, err := uc.GetCommentList(taskID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, commentsRes)
	}
}
