package middleware

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/delivery/response"
	"TaskFlow/internal/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthStateProvider interface {
	GetAuthState(ctx context.Context, id uuid.UUID) (domain.AuthState, error)
}

func JWTAuthMiddleware(cfg *config.Config, provider AuthStateProvider, logger *slog.Logger) gin.HandlerFunc {
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

		typ, ok := claims[domain.JWTType].(string)
		if !ok || typ != domain.JWTTypeAccess {
			abortUnauthorized(c, response.InvalidTokenClaims)
			return
		}

		id, ok := claims[domain.JWTID].(string)
		if !ok {
			abortUnauthorized(c, response.InvalidTokenClaims)
			return
		}

		idAsUUID, err := uuid.Parse(id)
		if err != nil {
			abortUnauthorized(c, response.InvalidTokenClaims)
			return
		}

		version, ok := claims[domain.JWTTokenVersion].(float64)
		if !ok {
			abortUnauthorized(c, response.InvalidTokenClaims)
			return
		}

		state, err := provider.GetAuthState(c.Request.Context(), idAsUUID)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				abortUnauthorized(c, response.InvalidToken)
				return
			}

			logger.Error("get auth state", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{Error: response.Internal})
			return
		}

		if int(version) != state.TokenVersion {
			abortUnauthorized(c, response.InvalidToken)
			return
		}

		c.Set(domain.JWTID, idAsUUID)
		c.Set(domain.CallerRole, state.Role)
	}
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
}
