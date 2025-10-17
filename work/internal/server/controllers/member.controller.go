package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/server/controllers/share"
	"idv/chris/MemoNest/internal/service"
)

type MemberController struct{}

// 使用者(註冊/登入)
func (u *MemberController) login(c *gin.Context) {
	account := "tester" // 暫用

	helper := share.NewSessionHelper(c)
	helper.Init(account)

	c.Redirect(http.StatusSeeOther, "/api/v1/article/list")
}

// 使用者登出
func (u *MemberController) logout(c *gin.Context) {
	// 使用 gin 取得 POST JSON 資料
	helper := share.NewSessionHelper(c)
	helper.Clear()
}

func NewMemberController(rg *gin.RouterGroup, di service.DI) {
	c := &MemberController{}
	r := rg.Group("/member")
	// r.POST("/login", c.login)
	// r.POST("/logout", c.logout)
	r.GET("/login", c.login)
	r.GET("/logout", c.logout)
}
