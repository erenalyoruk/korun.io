package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"korun.io/auth-service/internal/config"
	"korun.io/auth-service/internal/application"
)

func SetupRoutes(authService *application.AuthService, serverConfig *config.ServerConfig) *gin.Engine {
	router := gin.New()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = serverConfig.CORSAllowedOrigins
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(MetricsMiddleware())

	router.Use(RateLimiter())

	authHandler := NewAuthHandler(authService)

	router.GET("/health", authHandler.HealthCheck)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/:id/logout", authHandler.Logout)
		}
	}

	return router
}
