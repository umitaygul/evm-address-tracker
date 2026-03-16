package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/umitaygul/evm-address-tracker/internal/models"
	"github.com/umitaygul/evm-address-tracker/internal/repository"
)

type WebhookHandler struct {
	repo *repository.WebhookRepository
}

func NewWebhookHandler(repo *repository.WebhookRepository) *WebhookHandler {
	return &WebhookHandler{repo: repo}
}

func (h *WebhookHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		URL    string `json:"url" binding:"required"`
		Secret string `json:"secret"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	w, err := h.repo.Create(c.Request.Context(), userID, req.URL, req.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *WebhookHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	webhooks, err := h.repo.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if webhooks == nil {
		webhooks = []models.Webhook{}
	}
	c.JSON(http.StatusOK, webhooks)
}

func (h *WebhookHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
