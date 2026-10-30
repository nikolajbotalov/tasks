package comments

import (
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/middleware"
	commentsUC "TaskFlow/internal/use-cases/comment"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	auth gin.HandlerFunc,
	rateLimiter gin.HandlerFunc,
	uc commentsUC.UseCases,
	checker middleware.PermissionChecker,
	logger *slog.Logger,
) {
	commentRoutes := r.Group("/api/v1/tasks")
	commentRoutes.Use(rateLimiter)
	commentRoutes.Use(auth)
	{
		commentRoutes.GET("/:id/comments", middleware.RequirePermission(checker, domain.PermissionCommentRead),
			GetAllTaskComments(uc, logger))
		commentRoutes.POST("/:id/comments", middleware.RequirePermission(checker, domain.PermissionCommentCreate),
			CreateTaskComment(uc, logger))
	}
}
