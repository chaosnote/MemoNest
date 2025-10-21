package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/domain/service"
)

type MemberHandler struct {
	UC      *usecase.MemberUsecase
	Session service.Session
}

// 使用者(註冊/登入)
func (h *MemberHandler) Login(c *gin.Context) {
	account := "tester" // 暫用

	h.Session.Init(c)
	h.Session.SetAccount(account)

	c.Redirect(http.StatusSeeOther, "/api/v1/article/list")
}

// 使用者登出
func (h *MemberHandler) Logout(c *gin.Context) {
	h.Session.Init(c)
	h.Session.Clear()
}
