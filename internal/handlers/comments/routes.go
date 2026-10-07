package comments

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/handlers/middleware"
	commentsUC "TaskFlow/internal/use-cases/comment"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, uc commentsUC.UseCases) {
	commentRoutes := r.Group("/api/tasks")
	commentRoutes.Use(middleware.JWTAuthMiddleware(cfg))
	{
		commentRoutes.GET("/:id/comments", GetAllTaskComments(uc))
		commentRoutes.POST("/:id/comments", CreateTaskComment(uc))
	}
}
