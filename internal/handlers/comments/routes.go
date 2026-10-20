package comments

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/middleware"
	commentsUC "TaskFlow/internal/use-cases/comment"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	uc commentsUC.UseCases,
	checker middleware.PermissionChecker,
	logger *slog.Logger,
) {
	const (
		limit = 10
		burst = 20
	)

	commentRoutes := r.Group("/api/v1/tasks")
	commentRoutes.Use(middleware.JWTAuthMiddleware(cfg))
	{
		commentRoutes.GET("/:id/comments", middleware.RateLimiter(limit, burst),
			middleware.RequirePermission(checker, domain.PermissionCommentRead), GetAllTaskComments(uc, logger))
		commentRoutes.POST("/:id/comments", middleware.RateLimiter(limit, burst),
			middleware.RequirePermission(checker, domain.PermissionCommentCreate), CreateTaskComment(uc, logger))
	}
}
