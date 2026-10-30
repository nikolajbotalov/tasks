package middleware

import (
	"TaskFlow/internal/delivery/response"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	cleanupTTL      = 5 * time.Minute
	cleanupInterval = time.Minute
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func RateLimiter(limit rate.Limit, burst int) gin.HandlerFunc {
	var (
		mu          sync.Mutex
		clients     = make(map[string]*client)
		lastCleanup = time.Now()
	)

	removeStale := func(now time.Time) {
		if now.Sub(lastCleanup) < cleanupInterval {
			return
		}

		for ip, cl := range clients {
			if now.Sub(cl.lastSeen) > cleanupTTL {
				delete(clients, ip)
			}
		}

		lastCleanup = now
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		removeStale(now)

		cl, exists := clients[ip]
		if !exists {
			cl = &client{limiter: rate.NewLimiter(limit, burst)}
			clients[ip] = cl
		}
		cl.lastSeen = now
		mu.Unlock()

		reservation := cl.limiter.Reserve()
		delay := reservation.Delay()
		if delay == 0 {
			return
		}

		reservation.Cancel()

		c.Header("Retry-After", retryAfterHeader(delay))
		c.AbortWithStatusJSON(http.StatusTooManyRequests, response.ErrorResponse{
			Error: response.RateLimitExceeded,
		})
	}
}

func retryAfterHeader(delay time.Duration) string {
	seconds := int(math.Ceil(delay.Seconds()))
	if seconds < 1 {
		seconds = 1
	}

	return strconv.Itoa(seconds)
}
