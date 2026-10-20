package middleware

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	const bearerPrefix = "Bearer "

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": response.AuthorizationHeaderMissing})
			return
		}

		tokenString := authHeader[len(bearerPrefix):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			abortUnauthorized(c, response.InvalidToken)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			abortUnauthorized(c, response.InvalidTokenClaims)
			return
		}

		id, ok := claims[domain.JWTID].(string)
		if !ok {
			abortUnauthorized(c, response.InvalidTokenClaims)
			return
		}

		roleStr, ok := claims[domain.JWTRole].(string)
		if !ok {
			abortUnauthorized(c, response.InvalidTokenClaims)
			return
		}

		role := domain.Role(roleStr)
		if !role.IsValid() {
			abortUnauthorized(c, response.InvalidTokenRole)
			return
		}

		idAsUUID, err := uuid.Parse(id)
		if err != nil {
			abortUnauthorized(c, response.InvalidTokenClaims)
			return
		}

		c.Set(domain.JWTID, idAsUUID)
		c.Set(domain.JWTRole, role)
	}
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
}
