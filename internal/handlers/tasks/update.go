package tasks

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			log.Println("id is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		row := db.QueryRow("SELECT id FROM tasks WHERE id = $1", id)
		task := domain.Task{}

		if err := row.Scan(&task.ID); err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		if err := c.ShouldBindJSON(&task); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updatedAt := time.Now()

		_, err := db.Exec("UPDATE tasks SET (name, description, updated_at) = ($1, $2, $3) WHERE id = $4",
			task.Name, task.Description, updatedAt, id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "Task was updated"})
	}
}
