package middleware

import (
	"TaskFlow/internal/delivery/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	httpProtocol  = "http"
	httpsProtocol = "https"
	originHeader  = "Origin"
)

func RequireTrustedOrigin(allowedOrigins string) gin.HandlerFunc {
	allowed := make(map[string]struct{})

	for _, origin := range strings.Split(allowedOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader(originHeader)
		if origin == "" {
			c.Next()
			return
		}

		if origin == httpProtocol+c.Request.Host || origin == httpsProtocol+c.Request.Host {
			c.Next()
			return
		}

		if _, ok := allowed[origin]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse{Error: response.PermissionDenied})
			return
		}
	}
}
