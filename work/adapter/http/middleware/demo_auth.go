package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/service"
)

func DemoAuthMiddleware(deps service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !deps.Flag.AllowDemo {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "不允許使用 Demo 功能"})
			c.Abort()
			return
		}
		c.Next()
	}
}
