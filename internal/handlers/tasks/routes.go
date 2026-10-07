package tasks

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/handlers/middleware"
	tasksUC "TaskFlow/internal/use-cases/tasks"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine, cfg *config.Config, uc tasksUC.UseCases) {
	tasksProtected := g.Group("/api/tasks")
	tasksProtected.Use(middleware.JWTAuthMiddleware(cfg))
	{
		tasksProtected.GET("", GetAll(uc))
		tasksProtected.GET("/:id", GetByID(uc))
		tasksProtected.POST("", CreateTask(uc))
		tasksProtected.PUT("/:id", UpdateTask(uc))
		tasksProtected.DELETE("/:id", DeleteTask(uc))
	}
}
