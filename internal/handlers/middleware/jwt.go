package middleware

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenString := authHeader[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			abort(c, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			abort(c, "Invalid token claims")
			return
		}

		id, ok := claims[domain.JWTID].(string)
		if !ok {
			abort(c, "Invalid token claims")
			return
		}

		roleStr, ok := claims[domain.JWTRole].(string)
		if !ok {
			abort(c, "Invalid token claims")
			return
		}

		role := domain.Role(roleStr)
		if !role.IsValid() {
			abort(c, "Invalid token role")
			return
		}

		idAsUUID, err := uuid.Parse(id)
		if err != nil {
			abort(c, "Failed to parse uuid")
			return
		}

		c.Set(domain.JWTID, idAsUUID)
		c.Set(domain.JWTRole, role)
	}
}

func abort(c *gin.Context, errText string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errText})
}
