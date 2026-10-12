package router

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/handlers/auth"
	"TaskFlow/internal/handlers/comments"
	"TaskFlow/internal/handlers/middleware"
	"TaskFlow/internal/handlers/tasks"
	commentsUC "TaskFlow/internal/use-cases/comment"
	"net/http"

	"github.com/gin-gonic/gin"

	authUC "TaskFlow/internal/use-cases/auth"
	tasksUC "TaskFlow/internal/use-cases/tasks"
)

func AppRouters(
	authUc authUC.UseCases,
	cfg *config.Config,
	tasksUc tasksUC.UseCases,
	commentsUC commentsUC.UseCases,
	checker middleware.PermissionChecker,
) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	auth.SetupRoutes(r, authUc)
	tasks.SetupRoutes(r, cfg, tasksUc, checker)
	comments.SetupRoutes(r, cfg, commentsUC, checker)

	return r
}
