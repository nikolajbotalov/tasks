package auth

import (
	authUC "TaskFlow/internal/use-cases/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine, uc authUC.UseCases, logger *slog.Logger) {
	authRouters := g.Group("/api/v1/auth")
	{
		authRouters.POST("/register", RegisterUser(uc, logger))
		authRouters.POST("/login", LoginUser(uc, logger))
	}
}
