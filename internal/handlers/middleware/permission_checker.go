package middleware

import (
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/handlers/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionChecker interface {
	Has(role domain.Role, perm domain.Permission) bool
}

func RequirePermission(checker PermissionChecker, perm domain.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := helpers.GetRole(c)
		if err != nil || !checker.Has(role, perm) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": response.PermissionDenied})
			return
		}
	}
}
