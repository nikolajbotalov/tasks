package tasks

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteTask(db *sql.DB) gin.HandlerFunc {
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
			log.Println("task is not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "task is not found"})
			return
		}

		_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
		if err != nil {
			log.Println("failed to delete task")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Task was deleted"})
	}
}
