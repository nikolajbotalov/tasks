package tasks

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/handlers/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRouter(g *gin.Engine, db *sql.DB, cfg *config.Config) {
	tasksProtected := g.Group("/api/tasks")
	tasksProtected.Use(middleware.JWTAuthMiddleware(cfg))
	{
		tasksProtected.GET("", GetAll(db))
		tasksProtected.GET("/:id", GetByID(db))
		tasksProtected.POST("", CreateTask(db))
		tasksProtected.PUT("/:id", UpdateTask(db))
		tasksProtected.DELETE("/:id", DeleteTask(db))
	}
}
