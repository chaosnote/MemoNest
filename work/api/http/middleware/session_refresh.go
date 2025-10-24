package middleware

import (
	"fmt"
	"idv/chris/MemoNest/domain/service"

	"github.com/gin-gonic/gin"
)

func SessionRefresh(session service.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Update Session")
		session.Init(c)
		session.Refresh()
		c.Next()
	}
}
