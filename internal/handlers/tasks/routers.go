package tasks

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/handlers/middleware"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRouter(g *gin.Engine, db *sql.DB, cfg *config.Config, uc tasksUC.UseCases) {
	tasksProtected := g.Group("/api/tasks")
	tasksProtected.Use(middleware.JWTAuthMiddleware(cfg))
	{
		tasksProtected.GET("", GetAll(uc))
		tasksProtected.GET("/:id", GetByID(uc))
		tasksProtected.POST("", CreateTask(uc))
		tasksProtected.PUT("/:id", UpdateTask(db))
		tasksProtected.DELETE("/:id", DeleteTask(uc))
	}
}
