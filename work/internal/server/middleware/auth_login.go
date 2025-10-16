package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/service"
)

func MustLoginMiddleware(deps service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		flag, ok := s.Get(model.SK_IS_LOGIN).(bool)
		if !ok || !flag {
			fmt.Println("must login")
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}
