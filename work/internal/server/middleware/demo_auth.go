package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/service"
)

func DemoAuthMiddleware(deps service.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !deps.Flag.AllowDemo {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "不允許使用 Demo 功能"})
			return
		}
		c.Next()
	}
}
