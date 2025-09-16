package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jacoobjake/einvoice-api/config"
	"github.com/jacoobjake/einvoice-api/internal/handlers"
	"github.com/jacoobjake/einvoice-api/internal/repositories"
	"github.com/jacoobjake/einvoice-api/internal/services"
	"github.com/jacoobjake/einvoice-api/pkg/redisclient"
	"github.com/stephenafamo/bob"
)

func RegisterRoutes(r *gin.Engine, db bob.DB, cfg *config.Config, rdb *redisclient.RedisClient) {
	// Initialize repositories
	authTokenRepo := repositories.NewAuthTokenRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(authTokenRepo, userRepo, cfg, rdb)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Register routes
	apiGroup := r.Group("/api")
	{
		RegisterAuthRoutes(apiGroup, authHandler)
		// Add other route registrations here
	}
}
