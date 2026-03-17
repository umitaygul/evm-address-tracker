package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/umitaygul/evm-address-tracker/internal/api/handlers"
)

type mockUserRepo struct{}

func (m *mockUserRepo) Create(ctx interface{}, email, passwordHash string) (interface{}, error) {
	return nil, nil
}

func TestRegister_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	authHandler := handlers.NewAuthHandler(nil)
	r.POST("/register", authHandler.Register)

	body := bytes.NewBufferString(`{"email": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "error")
}

func TestLogin_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	authHandler := handlers.NewAuthHandler(nil)
	r.POST("/login", authHandler.Login)

	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "error")
}
