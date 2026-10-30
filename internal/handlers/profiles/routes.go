package profiles

import (
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/middleware"
	profilesUC "TaskFlow/internal/use-cases/profiles"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	auth gin.HandlerFunc,
	rateLimiter gin.HandlerFunc,
	uc profilesUC.UseCases,
	checker middleware.PermissionChecker,
	logger *slog.Logger,
) {
	profilesProtected := r.Group("/api/v1/profiles")
	profilesProtected.Use(rateLimiter)
	profilesProtected.Use(auth)
	{
		profilesProtected.GET("/:id", middleware.RequirePermission(checker, domain.PermissionProfileRead),
			GetProfileByUserID(uc, logger))
		profilesProtected.PUT("/:id", middleware.RequirePermission(checker, domain.PermissionProfileUpdate),
			UpdateProfile(uc, logger))
	}
}
