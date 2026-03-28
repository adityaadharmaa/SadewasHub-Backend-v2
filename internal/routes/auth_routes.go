package routes

import (
	"sadewashub-go/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.GET("/google", controllers.GoogleLogin)
		authGroup.GET("/google/callback", controllers.GoogleCallback)
		authGroup.POST("/login", controllers.Login)
	}
}
