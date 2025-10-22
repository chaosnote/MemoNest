package handle

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
)

type MemberHandler struct {
	UC      *usecase.MemberUsecase
	Session service.Session
}

// 使用者(註冊/登入)
func (h *MemberHandler) Login(c *gin.Context) {
	const msg = "login"
	logger := utils.NewFileLogger("./dist/logs/member/login", "console", 1)
	var err error
	defer func() {
		if err != nil {
			logger.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()
	var param struct {
		Account  string `json:"account" form:"account"`
		Password string `json:"password" form:"password"`
	}
	err = c.ShouldBind(&param)
	if err != nil {
		return
	}

	// 驗證帳號與密碼
	flag := h.UC.Login(param.Account, param.Password)
	if flag {
		h.Session.Init(c)
		h.Session.SetAccount(param.Account)
		c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
	} else {
		c.JSON(http.StatusOK, gin.H{"Code": "失敗原因", "message": ""})
	}
}

// 使用者登出
func (h *MemberHandler) Logout(c *gin.Context) {
	h.Session.Init(c)
	h.Session.Clear()
}

func (h *MemberHandler) Register(c *gin.Context) {
	time.Sleep(5 * time.Second)
	const msg = "register"
	logger := utils.NewFileLogger("./dist/logs/member/register", "console", 1)
	var err error
	defer func() {
		if err != nil {
			logger.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()
	c.JSON(http.StatusOK, gin.H{"Code": "測試", "message": ""})
}
