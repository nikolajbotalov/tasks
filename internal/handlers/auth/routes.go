package auth

import (
	"TaskFlow/internal/handlers"
	"TaskFlow/internal/handlers/middleware"
	authUC "TaskFlow/internal/use-cases/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine, uc authUC.UseCases, logger *slog.Logger) {
	authRouters := g.Group("/api/v1/auth")
	{
		authRouters.POST("/register", middleware.RateLimiter(handlers.UnprotectedRateLimit, handlers.UnprotectedRateBurst),
			RegisterUser(uc, logger))
		authRouters.POST("/login", middleware.RateLimiter(handlers.UnprotectedRateLimit, handlers.UnprotectedRateBurst),
			LoginUser(uc, logger))
	}
}
