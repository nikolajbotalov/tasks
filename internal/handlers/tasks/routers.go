package tasks

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRouter(g *gin.Engine, db *sql.DB) {
	taskRouters := g.Group("/api/tasks")
	{
		taskRouters.GET("/", GetAll(db))
		taskRouters.GET("/:id", GetByID(db))
		taskRouters.POST("/", CreateTask(db))
		taskRouters.PUT("/:id", UpdateTask(db))
		taskRouters.DELETE("/:id", DeleteTask(db))
	}
}
