package tasks

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		row := db.QueryRow("SELECT * FROM tasks WHERE id=$1", id)
		t := domain.Task{}

		if err := row.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, t)
	}
}
