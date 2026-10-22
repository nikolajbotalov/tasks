package roles

import (
	"TaskFlow/internal/delivery/response"
	rolesUC "TaskFlow/internal/use-cases/roles"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllUsersRoles(uc rolesUC.UseCases, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := uc.GetAllUsersRoles()
		if err != nil {
			logger.Error("get user roles", slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": response.Internal})
			return
		}

		c.JSON(http.StatusOK, roles)
	}
}
