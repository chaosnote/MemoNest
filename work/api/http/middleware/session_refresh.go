package middleware

import (
	"idv/chris/MemoNest/domain/service"

	"github.com/gin-gonic/gin"
)

func SessionRefresh(session service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.Init(c)
		session.Refresh()
		c.Next()
	}
}
