package http

import (
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	SessionKeyAccount = "account"
	SessionKeyAESKey  = "aes_key"
	SessionKeyIsLogin = "is_login"
	DefaultAccount    = "guest"
	DefaultAESKey     = "default-key"
)

type GinSession struct {
	store sessions.Session
}

func NewGinSession(ctx *gin.Context) service.Session {
	return &GinSession{store: sessions.Default(ctx)}
}

func (s *GinSession) Init(account string) {
	s.store.Set(SessionKeyAccount, account)
	s.store.Set(SessionKeyAESKey, utils.MD5Hash(account)) // 可改為 row_id
	s.store.Set(SessionKeyIsLogin, true)
	s.store.Save()
}

func (s *GinSession) Clear() {
	s.store.Clear()
	s.store.Save()
}

func (s *GinSession) GetAccount() string {
	val := s.store.Get(SessionKeyAccount)
	if str, ok := val.(string); ok {
		return str
	}
	return DefaultAccount
}

func (s *GinSession) GetAESKey() string {
	val := s.store.Get(SessionKeyAESKey)
	if str, ok := val.(string); ok {
		return str
	}
	return DefaultAESKey
}

func (s *GinSession) IsLogin() bool {
	val := s.store.Get(SessionKeyIsLogin)
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}
