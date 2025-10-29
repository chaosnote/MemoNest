package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/utils"
)

type MemberHandler struct {
	CommonHandler

	UC *usecase.MemberUsecase
}

// 使用者(註冊/登入)
func (h *MemberHandler) Login(c *gin.Context) {
	const msg = "login"

	var err error
	defer func() {
		if err != nil {
			h.Log.Error(msg, zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"Code": err.Error()})
		}
	}()

	var param struct {
		Account  string `json:"account" form:"account"`
		Password string `json:"password" form:"password"`
		Token    string `json:"token" form:"token"`
	}
	err = c.ShouldBind(&param)
	if err != nil {
		return
	}

	aes_key := utils.MD5Hash(c.ClientIP())
	if len(param.Token) > 0 {
		tmp, err := utils.AesDecrypt(param.Token, []byte(aes_key))
		if err != nil {
			return
		}
		list := strings.Split(tmp, "|")
		if len(list) != 2 {
			c.JSON(http.StatusOK, gin.H{"Code": "請使用手動登入"})
			return
		}
		param.Account = list[0]
		param.Password = list[1]
	}

	h.Log.Info(msg, zap.Any("params", param))

	_, err = h.UC.Login(param.Account, param.Password, c.ClientIP())
	if err != nil {
		return
	}

	token, err := utils.AesEncrypt([]byte(fmt.Sprintf("%s|%s", param.Account, param.Password)), []byte(aes_key))
	if err != nil {
		return
	}

	h.Session.Init(c)
	h.Session.SetAccount(param.Account, c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"Code":    "OK",
		"message": "",
		"token":   token,
	})
}

// 使用者登出
func (h *MemberHandler) Logout(c *gin.Context) {
	h.Session.Init(c)
	h.Session.Clear()

	c.Redirect(http.StatusFound, "/")
}

func (h *MemberHandler) Register(c *gin.Context) {
	const msg = "register"

	var err error
	defer func() {
		if err != nil {
			h.Log.Error(msg, zap.Error(err))
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

	h.Log.Info(msg, zap.Any("params", param))

	_, err = h.UC.Register(param.Account, param.Password, c.ClientIP())
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"Code": "OK", "message": ""})
}
