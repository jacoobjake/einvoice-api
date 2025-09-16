package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jacoobjake/einvoice-api/internal/services"
	"github.com/jacoobjake/einvoice-api/pkg/response"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, response.JSONApiResponse{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := authService.VerifyToken(c.Request.Context(), token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, response.JSONApiResponse{
				Success: false,
				Message: "Invalid token",
			})
			c.Abort()
			return
		}

		// Set authorized user in context
		c.Set("user", user)
		c.Set("auth_token", token)

		c.Next()
	}
}
