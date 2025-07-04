package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"korun.io/secret-service/internal/service"
	"korun.io/shared/models"
)

type SecretHandler struct {
	secretService *service.SecretService
}

var (
	ErrIdParameterRequired   = "id parameter is required"
	ErrNameParameterRequired = "name parameter is required"
	ErrInvalidRequest        = "invalid request"
	ErrInvalidSecretID       = "invalid secret id"
	ErrFailedToRetrieve      = "failed to retrieve secret"
	ErrSecretNotFound        = "secret not found"
)

func NewSecretHandler(secretService *service.SecretService) *SecretHandler {
	return &SecretHandler{secretService: secretService}
}

func (h *SecretHandler) GetSecrets(ctx *gin.Context) {
	secrets, err := h.secretService.GetSecrets(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToRetrieve})
		return
	}

	ctx.JSON(http.StatusOK, secrets)
}

func (h *SecretHandler) GetSecretByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidSecretID})
		return
	}

	secret, err := h.secretService.GetSecretByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToRetrieve})
		return
	}

	if secret == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": ErrSecretNotFound})
		return
	}

	ctx.JSON(http.StatusOK, secret)
}

func (h *SecretHandler) GetSecretByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrNameParameterRequired})
		return
	}

	secret, err := h.secretService.GetSecretByName(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToRetrieve})
		return
	}

	if secret == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": ErrSecretNotFound})
		return
	}

	ctx.JSON(http.StatusOK, secret)
}

func (h *SecretHandler) CreateSecret(ctx *gin.Context) {
	var secret models.Secret
	if err := ctx.ShouldBindJSON(&secret); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}

	if err := h.secretService.CreateSecret(ctx, &secret); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToRetrieve})
		return
	}

	ctx.JSON(http.StatusCreated, secret)
}

func (h *SecretHandler) UpdateSecret(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidSecretID})
		return
	}

	var secret models.Secret
	if err := ctx.ShouldBindJSON(&secret); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}

	secret.ID = uint(id)
	if err := h.secretService.UpdateSecret(ctx, &secret); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToRetrieve})
		return
	}

	ctx.JSON(http.StatusOK, secret)
}

func (h *SecretHandler) DeleteSecret(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidSecretID})
		return
	}

	if err := h.secretService.DeleteSecret(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrSecretNotFound})
		return
	}

	ctx.Status(http.StatusNoContent)
}
