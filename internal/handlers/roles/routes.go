package roles

import (
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/middleware"
	"log/slog"

	rolesUC "TaskFlow/internal/use-cases/roles"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	auth gin.HandlerFunc,
	rateLimiter gin.HandlerFunc,
	uc rolesUC.UseCases,
	checker middleware.PermissionChecker,
	logger *slog.Logger,
) {
	usersProtected := r.Group("/api/v1/users")
	usersProtected.Use(rateLimiter)
	usersProtected.Use(auth)
	{
		usersProtected.GET("/roles", middleware.RequirePermission(checker, domain.PermissionUserRead),
			GetAllUsersRoles(uc, logger))
		usersProtected.PUT("/:id/roles", middleware.RequirePermission(checker, domain.PermissionUserUpdate),
			UpdateUserRole(uc, logger))
	}
}
