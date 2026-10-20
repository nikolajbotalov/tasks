package tasks

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/middleware"
	tasksUC "TaskFlow/internal/use-cases/tasks"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	uc tasksUC.UseCases,
	checker middleware.PermissionChecker,
	logger *slog.Logger,
) {
	const (
		limit = 10
		burst = 20
	)

	tasksProtected := r.Group("/api/v1/tasks")
	tasksProtected.Use(middleware.JWTAuthMiddleware(cfg))
	{
		tasksProtected.GET("", middleware.RateLimiter(limit, burst),
			middleware.RequirePermission(checker, domain.PermissionTaskRead), GetAll(uc, logger))
		tasksProtected.GET("/:id", middleware.RateLimiter(limit, burst),
			middleware.RequirePermission(checker, domain.PermissionTaskRead), GetByID(uc, logger))
		tasksProtected.POST("", middleware.RateLimiter(limit, burst),
			middleware.RequirePermission(checker, domain.PermissionTaskCreate), CreateTask(uc, logger))
		tasksProtected.PUT("/:id", middleware.RateLimiter(limit, burst),
			middleware.RequirePermission(checker, domain.PermissionTaskUpdate), UpdateTask(uc, logger))
		tasksProtected.DELETE("/:id", middleware.RateLimiter(limit, burst),
			middleware.RequirePermission(checker, domain.PermissionTaskDelete), DeleteTask(uc, logger))
	}
}
