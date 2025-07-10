package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	tokenHandler := NewTokenHandler()

	router.GET("/health", tokenHandler.HealthCheck)

	return router
}
