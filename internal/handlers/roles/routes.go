package roles

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers"
	"TaskFlow/internal/handlers/middleware"
	"log/slog"

	rolesUC "TaskFlow/internal/use-cases/roles"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	uc rolesUC.UseCases,
	checker middleware.PermissionChecker,
	logger *slog.Logger,
) {
	usersProtected := r.Group("/api/v1/users")
	usersProtected.Use(middleware.JWTAuthMiddleware(cfg))
	{
		usersProtected.GET("/roles", middleware.RateLimiter(handlers.ProtectedRateLimit, handlers.ProtectedRateBurst),
			middleware.RequirePermission(checker, domain.PermissionUserRead), GetAllUsersRoles(uc, logger))
		usersProtected.PUT("/:id/roles", middleware.RateLimiter(handlers.ProtectedRateLimit, handlers.ProtectedRateBurst),
			middleware.RequirePermission(checker, domain.PermissionUserUpdate), UpdateUserRole(uc, logger))
	}
}
