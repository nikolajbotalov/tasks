package auth

import (
	"TaskFlow/internal/handlers/middleware"
	authUC "TaskFlow/internal/use-cases/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine, uc authUC.UseCases, logger *slog.Logger) {
	const (
		limit = 0.1
		burst = 5
	)

	authRouters := g.Group("/api/v1/auth")
	{
		authRouters.POST("/register", middleware.RateLimiter(limit, burst), RegisterUser(uc, logger))
		authRouters.POST("/login", middleware.RateLimiter(limit, burst), LoginUser(uc, logger))
	}
}
