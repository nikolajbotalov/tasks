package tasks

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAll(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")
		rows, err := db.Query("SELECT * FROM tasks LIMIT $1 OFFSET $2", limit, page)
		if err != nil {
			log.Println(err)
		}

		tasks := make([]domain.Task, 0)

		for rows.Next() {
			t := domain.Task{}
			err = rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.UpdatedAt)
			if err != nil {
				log.Println(err)
				continue
			}
			tasks = append(tasks, t)
		}

		log.Println(tasks)

		limitNumber, err := strconv.Atoi(limit)
		if err != nil {
			log.Println(err)
		}

		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, domain.GetAllTasksResponse{
			Tasks: tasks,
			Limit: limitNumber,
			Page:  pageNumber,
			Total: len(tasks),
		})
	}
}
