package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	xxx "idv/chris/MemoNest/adapter/http"
	"idv/chris/MemoNest/service"
)

type MemberController struct{}

// 使用者(註冊/登入)
func (u *MemberController) login(c *gin.Context) {
	account := "tester" // 暫用

	helper := xxx.NewGinSession(c)
	helper.Init(account)

	c.Redirect(http.StatusSeeOther, "/api/v1/article/list")
}

// 使用者登出
func (u *MemberController) logout(c *gin.Context) {
	// 使用 gin 取得 POST JSON 資料
	helper := xxx.NewGinSession(c)
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
