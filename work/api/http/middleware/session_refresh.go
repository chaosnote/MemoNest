package middleware

import (
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/domain/service"
)

func SessionRefresh(session service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.Init(c)
		session.Refresh()
		c.Next()
	}
}
