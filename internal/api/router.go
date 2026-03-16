package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/umitaygul/evm-address-tracker/internal/api/handlers"
	"github.com/umitaygul/evm-address-tracker/internal/api/middleware"
	"github.com/umitaygul/evm-address-tracker/internal/repository"
)

func NewRouter(db *pgxpool.Pool) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Repositories
	userRepo := repository.NewUserRepository(db)
	addressRepo := repository.NewAddressRepository(db)

	// Handlers
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(userRepo)
	addressHandler := handlers.NewAddressHandler(addressRepo)

	// Public routes
	r.GET("/health", healthHandler.GetHealth)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.Auth())
		{
			addresses := protected.Group("/addresses")
			{
				addresses.POST("", addressHandler.Add)
				addresses.GET("", addressHandler.List)
				addresses.GET("/:id", addressHandler.Get)
				addresses.DELETE("/:id", addressHandler.Delete)
			}
		}
	}

	return r
}
