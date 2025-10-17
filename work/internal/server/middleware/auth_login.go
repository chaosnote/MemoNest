package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/server/controllers/share"
	"idv/chris/MemoNest/internal/service"
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
