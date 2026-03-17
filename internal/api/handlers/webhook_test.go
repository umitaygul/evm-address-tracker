package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/umitaygul/evm-address-tracker/internal/api/handlers"
)

func TestWebhook_MissingUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	webhookHandler := handlers.NewWebhookHandler(nil)
	r.POST("/webhooks", webhookHandler.Create)

	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/webhooks", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWebhook_InvalidJson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	webhookHandler := handlers.NewWebhookHandler(nil)
	r.POST("/webhooks", webhookHandler.Create)

	body := bytes.NewBufferString(`{bozuk json`)
	req := httptest.NewRequest(http.MethodPost, "/webhooks", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestWebhook_ValidRequest_ReachesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	webhookHandler := handlers.NewWebhookHandler(nil)
	r.POST("/webhooks", webhookHandler.Create)

	body := bytes.NewBufferString(`{"url": "https://example.com", "secret": "mysecret"}`)
	req := httptest.NewRequest(http.MethodPost, "/webhooks", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// DB nil olduğu için 500 döner ama 400 veya 401 dönmedi
	// yani validation geçti, handler'a ulaştı
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
	assert.NotEqual(t, http.StatusUnauthorized, w.Code)
}
