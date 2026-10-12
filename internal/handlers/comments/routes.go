package comments

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/middleware"
	commentsUC "TaskFlow/internal/use-cases/comment"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, uc commentsUC.UseCases, checker middleware.PermissionChecker) {
	commentRoutes := r.Group("/api/tasks")
	commentRoutes.Use(middleware.JWTAuthMiddleware(cfg))
	{
		commentRoutes.GET("/:id/comments", middleware.RequirePermission(checker, domain.PermissionCommentRead), GetAllTaskComments(uc))
		commentRoutes.POST("/:id/comments", middleware.RequirePermission(checker, domain.PermissionCommentCreate), CreateTaskComment(uc))
	}
}
