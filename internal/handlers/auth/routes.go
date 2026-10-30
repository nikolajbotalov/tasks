package auth

import (
	authUC "TaskFlow/internal/use-cases/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine, rateLimiter gin.HandlerFunc, uc authUC.UseCases, logger *slog.Logger) {
	authRouters := g.Group("/api/v1/auth")
	authRouters.Use(rateLimiter)
	{
		authRouters.POST("/register", RegisterUser(uc, logger))
		authRouters.POST("/login", LoginUser(uc, logger))
	}
}
