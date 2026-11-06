package auth

import (
	"TaskFlow/internal/config"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const refreshCookieName = "taskflow_refresh_token"

func setRefreshCookie(c *gin.Context, cfg *config.Config, token string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	if maxAge <= 0 {
		maxAge = 1
	}

	c.SetCookieData(&http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     cfg.Cookie.Path,
		Domain:   cfg.Cookie.Domain,
		MaxAge:   maxAge,
		Secure:   cfg.Cookie.Secure,
		HttpOnly: true,
		SameSite: cookieSameSite(cfg.Cookie.SameSite),
	})
}

func clearRefreshCookie(c *gin.Context, cfg *config.Config) {
	c.SetCookieData(&http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     cfg.Cookie.Path,
		Domain:   cfg.Cookie.Domain,
		MaxAge:   -1,
		Secure:   cfg.Cookie.Secure,
		HttpOnly: true,
		SameSite: cookieSameSite(cfg.Cookie.SameSite),
	})
}

func cookieSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
