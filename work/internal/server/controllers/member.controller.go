package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/service"
)

type MemberController struct{}

// 使用者(註冊/登入)
func (u *MemberController) login(c *gin.Context) {
	s := sessions.Default(c)
	s.Set(model.SK_ACCOUNT, "tester")
	s.Set(model.SK_IS_LOGIN, true)
	s.Save()

	c.Redirect(http.StatusSeeOther, "/api/v1/article/fresh")
}

// 使用者登出
func (u *MemberController) logout(c *gin.Context) {
	// 使用 gin 取得 POST JSON 資料
	s := sessions.Default(c)
	s.Clear()
	s.Save()
}

func NewMemberController(rg *gin.RouterGroup, di service.DI) {
	c := &MemberController{}
	r := rg.Group("/member")
	// r.POST("/login", c.login)
	// r.POST("/logout", c.logout)
	r.GET("/login", c.login)
	r.GET("/logout", c.logout)
}
