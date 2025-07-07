package api

import (
	"github.com/gin-gonic/gin"
	"korun.io/auth-service/internal/service"
)

func SetupRoutes(authService *service.AuthService) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	authHandler := NewAuthHandler(authService)

	router.GET("/health", authHandler.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}
	}

	return router
}
