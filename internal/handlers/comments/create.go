package comments

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"TaskFlow/internal/handlers/validation"
	commentsUC "TaskFlow/internal/use-cases/comment"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCommentRequest struct {
	Text string `json:"text" binding:"required,min=2,max=10000"`
}

func CreateTaskComment(uc commentsUC.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, authorID, err := helpers.GetTaskAndAuthorIDs(c)
		if err != nil {
			if errors.Is(err, domain.ErrGetTaskID) {
				c.JSON(http.StatusBadRequest, gin.H{"error": response.IncorrectTaskID})
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": response.IncorrectAuthorID})
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
				c.JSON(http.StatusNotFound, gin.H{"error": response.TaskNotFound})
				return
			}

			fmt.Printf("create comment for task %s, reason: %v\n", taskID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusCreated, taskComment)
	}
}
