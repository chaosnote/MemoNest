package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/service"
)

type MemberController struct{}

// 使用者(註冊/登入)
func (u *MemberController) login(service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
}

// 使用者登出
func (u *MemberController) logout(service.DI) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 gin 取得 POST JSON 資料
	}
}

func NewMemberController(rg *gin.RouterGroup, di service.DI) {
	c := &MemberController{}
	r := rg.Group("/member")
	r.POST("/login", c.login(di))
	r.POST("/logout", c.logout(di))
}
