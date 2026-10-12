package tasks

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/middleware"
	tasksUC "TaskFlow/internal/use-cases/tasks"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine, cfg *config.Config, uc tasksUC.UseCases, checker middleware.PermissionChecker) {
	tasksProtected := g.Group("/api/tasks")
	tasksProtected.Use(middleware.JWTAuthMiddleware(cfg))
	{
		tasksProtected.GET("", middleware.RequirePermission(checker, domain.PermissionTaskRead), GetAll(uc))
		tasksProtected.GET("/:id", middleware.RequirePermission(checker, domain.PermissionTaskRead), GetByID(uc))
		tasksProtected.POST("", middleware.RequirePermission(checker, domain.PermissionTaskCreate), CreateTask(uc))
		tasksProtected.PUT("/:id", middleware.RequirePermission(checker, domain.PermissionTaskUpdate), UpdateTask(uc))
		tasksProtected.DELETE("/:id", middleware.RequirePermission(checker, domain.PermissionTaskDelete), DeleteTask(uc))
	}
}
