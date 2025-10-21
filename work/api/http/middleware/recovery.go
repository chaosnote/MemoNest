package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinRecovery 攔截 panic 並且記錄，回 500
func GinRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", zap.Any("error", r))
				c.AbortWithStatus(http.StatusInternalServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
