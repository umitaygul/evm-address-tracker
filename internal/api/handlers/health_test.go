package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/umitaygul/evm-address-tracker/internal/api/handlers"
)

func TestGetHealth_DBDown(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := handlers.NewHealthHandler(nil)

	r := gin.New()
	r.GET("/health", handler.GetHealth)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "unhealthy", body["status"])
	assert.Equal(t, "down", body["db"])
}
