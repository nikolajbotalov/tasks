package helpers

import (
	"TaskFlow/internal/domain"

	"github.com/gin-gonic/gin"
)

func GetRole(c *gin.Context) (domain.Role, error) {
	roleParam, exists := c.Get(domain.JWTRole)
	if !exists {
		return "", domain.ErrUnknownRole
	}

	role, ok := roleParam.(domain.Role)
	if !ok {
		return "", domain.ErrUnknownRole
	}

	return role, nil
}
