package api

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RequestIDKey = "request_id"

func RequestIDMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.NewString()

		c.Set(RequestIDKey, reqID)

		requestLogger := logger.With(zap.String("request_id", reqID))

		c.Set("logger", requestLogger)

		c.Next()
	}
}

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqLogger := getLoggerFromContext(c)

		defer func() {
			if err := recover(); err != nil {
				reqLogger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
					zap.String("request_uri", c.Request.RequestURI),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "An internal server error occurred.",
				})
			}
		}()
		c.Next()
	}
}

func ZapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		reqLogger := getLoggerFromContext(c)

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		fields := []zap.Field{
			zap.Duration("latency", latency),
			zap.Int("status_code", statusCode),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if statusCode >= 400 {
			reqLogger.Warn("Request completed with error", append(fields, zap.Int("status_code", statusCode))...)
		} else {
			reqLogger.Info("Request completed successfully", append(fields, zap.Int("status_code", statusCode))...)
		}
	}
}

func getLoggerFromContext(c *gin.Context) *zap.Logger {
	if logger, ok := c.Get("logger"); ok {
		if reqLogger, ok := logger.(*zap.Logger); ok {
			return reqLogger
		}
	}

	return zap.NewNop()
}
