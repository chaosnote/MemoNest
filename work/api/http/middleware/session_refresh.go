package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/domain/service"
)

func SessionRefresh(session service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.Init(c)
		err := session.Refresh()
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
