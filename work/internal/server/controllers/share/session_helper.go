package share

import (
	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type session_store struct {
	s sessions.Session
}

func (ss session_store) Init(account string) {
	ss.s.Set(model.SK_ACCOUNT, account)
	ss.s.Set(model.SK_AES_KEY, utils.MD5Hash(account))
	ss.s.Set(model.SK_IS_LOGIN, true)
	ss.s.Save()
}

func (ss session_store) Clear() {
	ss.s.Clear()
	ss.s.Save()
}

func (ss session_store) GetAESKey() string {
	return ss.s.Get(model.SK_AES_KEY).(string)
}

func (ss session_store) GetAccount() string {
	return ss.s.Get(model.SK_ACCOUNT).(string)
}

func (ss session_store) IsLogin() bool {
	flag, ok := ss.s.Get(model.SK_IS_LOGIN).(bool)
	if !ok {
		return false
	}
	return flag
}

//-----------------------------------------------

type SessionImpl interface {
	Init(account string)
	Clear()

	GetAESKey() string  // 加密 Key
	GetAccount() string // 使用者帳號
	IsLogin() bool      // {true:已登入, false:未登入}
}

func NewSessionHelper(c *gin.Context) SessionImpl {
	return &session_store{
		s: sessions.Default(c),
	}
}
