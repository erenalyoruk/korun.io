package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRoutes(logger *zap.Logger) *gin.Engine {
	router := gin.New()

	router.Use(RequestIDMiddleware(logger))
	router.Use(ZapLoggerMiddleware())

	tokenHandler := NewTokenHandler()

	router.GET("/health", tokenHandler.HealthCheck)

	return router
}
