package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/utils"
)

// IPCheckMiddleware 是用來檢查 IP 是否為私有 IP 的中間件
func IPCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 context 取得請求的遠端 IP
		clientIP := net.ParseIP(c.RemoteIP())

		// 檢查 IP 是否為私有 IP，如果不是，就中止請求並回傳 403 錯誤
		if !utils.IsPrivateIP(clientIP) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "只允許區域 IP",
			})
			c.Abort() // 停止執行後續的處理函式
			return
		}

		// 如果 IP 檢查通過，繼續執行後續的處理函式
		c.Next()
	}
}
