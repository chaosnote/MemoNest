package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/server/controllers/share"
	"idv/chris/MemoNest/service"
)

func MustLoginMiddleware(deps service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		helper := share.NewSessionHelper(c)
		if !helper.IsLogin() {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}
