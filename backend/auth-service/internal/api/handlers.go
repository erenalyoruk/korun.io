package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"korun.io/auth-service/internal/service"
	"korun.io/shared/models"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAccountExists):
			c.JSON(http.StatusConflict, gin.H{"error": "account already exists"})
		case errors.Is(err, service.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, service.ErrPasswordTooWeak):
			c.JSON(http.StatusBadRequest, gin.H{"error": "password is too weak"})
		default:
			slog.Error("Registration failed", "error", err, "email", req.Email)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register account"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	response, err := h.authService.Login(c.Request.Context(), &req, clientIP, userAgent)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		slog.Error("Login failed", "error", err, "email", req.Email, "ip", clientIP)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	response, err := h.authService.RefreshToken(c.Request.Context(), &req, clientIP, userAgent)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}
		slog.Error("Token refresh failed", "error", err, "ip", clientIP)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), accountID); err != nil {
		slog.Error("Logout failed", "error", err, "account_id", accountID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "auth-service",
	})
}
