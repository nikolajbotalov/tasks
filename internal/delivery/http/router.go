package router

import (
	"TaskFlow/internal/handlers/auth"
	"TaskFlow/internal/handlers/tasks"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	authUC "TaskFlow/internal/use-cases/auth"
)

func AppRouters(db *sql.DB, uc authUC.UseCases) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	tasks.SetupRouter(r, db)
	auth.SetupRouter(r, uc)

	return r
}
