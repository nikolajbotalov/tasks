package auth

import (
	authUC "TaskFlow/internal/use-cases/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine, uc authUC.UseCases) {
	authRouters := g.Group("/api/auth")
	{
		authRouters.POST("/register", RegisterUser(uc))
		authRouters.POST("/login", LoginUser(uc))
	}
}
