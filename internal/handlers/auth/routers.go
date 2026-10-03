package auth

import (
	"TaskFlow/internal/use-cases/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine, uc auth.UseCases) {
	authRouters := g.Group("/api/auth")
	{
		authRouters.POST("/register", RegisterUser(uc))
		authRouters.POST("/login", LoginUser(uc))
	}
}
