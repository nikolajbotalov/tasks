package router

import (
	"TaskFlow/internal/handlers/tasks"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppRouters(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	tasks.SetupRouter(r, db)

	return r
}
