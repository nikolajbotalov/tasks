package users

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine, db *sql.DB) {
	usersRouters := g.Group("/api/tasks")
	{

	}
}
