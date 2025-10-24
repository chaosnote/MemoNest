package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/domain/service"
)

func Auth(session service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.Init(c)
		if !session.IsLogin() || session.GetIP() != c.ClientIP() {
			session.Clear()
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}
