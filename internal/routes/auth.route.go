package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jacoobjake/einvoice-api/internal/handlers"
	"github.com/jacoobjake/einvoice-api/internal/routes/middlewares"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, handler *handlers.AuthHandler) {

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", handler.Login)

		authGroup.Use(middlewares.AuthMiddleware(handler.AuthService))
		authGroup.POST("/logout", handler.Logout)
	}
}
