package tasks

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/handlers/helpers"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAll(uc tasksUC.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID, err := helpers.GetAuthorID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.IncorrectAuthorID)
			return
		}

		page := c.DefaultQuery("page", "1")
		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			fmt.Printf("parsing %s, invalid syntax: %v\n", page, err)
			c.JSON(http.StatusBadRequest, response.InvalidPageParam)
			return
		}
		if pageNumber < 1 {
			c.JSON(http.StatusBadRequest, response.InvalidPageParam)
			return
		}

		limit := c.DefaultQuery("limit", "10")
		limitNumber, err := strconv.Atoi(limit)
		if err != nil {
			fmt.Printf("parsing %s, invalid syntax: %v\n", limit, err)
			c.JSON(http.StatusBadRequest, response.InvalidLimitParam)
			return
		}
		if limitNumber <= 0 || limitNumber > 100 {
			c.JSON(http.StatusBadRequest, response.InvalidLimitParam)
			return
		}

		tasks, err := uc.GetAllTasks(authorID, pageNumber, limitNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
