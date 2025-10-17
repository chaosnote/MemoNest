package share

import (
	"idv/chris/MemoNest/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	sk_account  string = "account"
	sk_aes_key  string = "aes_key"
	sk_is_login string = "is_login"
)

//-----------------------------------------------

type session_store struct {
	s sessions.Session
}

func (ss session_store) Init(account string) {
	ss.s.Set(sk_account, account)
	ss.s.Set(sk_aes_key, utils.MD5Hash(account))
	ss.s.Set(sk_is_login, true)
	ss.s.Save()
}

func (ss session_store) Clear() {
	ss.s.Clear()
	ss.s.Save()
}

func (ss session_store) GetAESKey() string {
	return ss.s.Get(sk_aes_key).(string)
}

func (ss session_store) GetAccount() string {
	return ss.s.Get(sk_account).(string)
}

func (ss session_store) IsLogin() bool {
	flag, ok := ss.s.Get(sk_is_login).(bool)
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
