package middleware

import (
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/domain/service"
)

func RecordPath(session service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.Init(c)
		session.SetURL(c.Request.URL.Path)
		c.Next()
	}
}
