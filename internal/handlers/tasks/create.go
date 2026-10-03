package tasks

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		newUUID := uuid.New()
		createdAt := time.Now()
		updatedAt := time.Now()

		var task domain.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO tasks (id, name, description, created_at, updated_at) values ($1, $2, $3, $4, $5)",
			newUUID.String(), task.Name, task.Description, createdAt, updatedAt)
		if err != nil {
			log.Println(err)
		}

		row := db.QueryRow("SELECT * FROM tasks WHERE  id = $1", newUUID)

		result := domain.Task{}

		err = row.Scan(&result.ID, &result.Name, &result.Description, &result.CreatedAt, &result.UpdatedAt)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusCreated, result)
	}
}
