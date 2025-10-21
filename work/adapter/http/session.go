package http

import (
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	SK_Account = "account"
	SK_AESKey  = "aes_key"
	SK_IsLogin = "is_login"
)

type GinSession struct {
	store sessions.Session
}

func (s *GinSession) Init(ctx *gin.Context) {
	s.store = sessions.Default(ctx)
}

func (s *GinSession) Clear() {
	s.store.Clear()
	s.store.Save()
}

func (s *GinSession) SetAccount(account string) {
	s.store.Set(SK_Account, account)
	s.store.Set(SK_AESKey, utils.MD5Hash(account)) // 可改為 row_id
	s.store.Set(SK_IsLogin, true)
	s.store.Save()
}

func (s *GinSession) GetAccount() string {
	val := s.store.Get(SK_Account)
	str, ok := val.(string)
	if ok {
		return str
	}
	return ""
}

func (s *GinSession) GetAESKey() string {
	val := s.store.Get(SK_AESKey)
	str, ok := val.(string)
	if ok {
		return str
	}
	return ""
}

func (s *GinSession) IsLogin() bool {
	val := s.store.Get(SK_IsLogin)
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

func (s *GinSession) Refresh() {
	s.store.Save()
}

//-----------------------------------------------

func NewGinSession() service.Session {
	return &GinSession{}
}
