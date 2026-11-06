package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/handlers/middleware"
	authUC "TaskFlow/internal/use-cases/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	g *gin.Engine,
	cfg *config.Config,
	strictLimiter gin.HandlerFunc,
	refreshLimiter gin.HandlerFunc,
	uc authUC.UseCases,
	logger *slog.Logger,
) {
	authRouters := g.Group("/api/v1/auth")
	authRouters.Use(middleware.RequireTrustedOrigin(cfg.CORS.AllowedOrigins))
	{
		authRouters.POST("/register", strictLimiter, RegisterUser(uc, logger))
		authRouters.POST("/login", strictLimiter, LoginUser(uc, logger, cfg))
		authRouters.POST("/refresh", refreshLimiter, RefreshToken(uc, logger, cfg))
		authRouters.POST("/logout", refreshLimiter, Logout(uc, logger, cfg))
	}
}
