package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umitaygul/evm-address-tracker/internal/api/handlers"
)

func NewRouter(db *pgxpool.Pool) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// Handlers
	healthHandler := handlers.NewHealthHandler(db)

	// Routes
	r.GET("/health", healthHandler.GetHealth)

	return r
}
