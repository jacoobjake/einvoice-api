package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jacoobjake/einvoice-api/pkg"
)

func ClientIpMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipStr := c.ClientIP()

		ctx := pkg.SetCtxClientIp(c.Request.Context(), ipStr)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
