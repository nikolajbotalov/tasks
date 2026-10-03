package router

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/handlers/auth"
	"TaskFlow/internal/handlers/tasks"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	authUC "TaskFlow/internal/use-cases/auth"
)

func AppRouters(db *sql.DB, uc authUC.UseCases, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	auth.SetupRouter(r, uc)
	tasks.SetupRouter(r, db, cfg)

	return r
}
